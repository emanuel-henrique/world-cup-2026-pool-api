// internal/matches/handler.go
package matches

import (
	"fmt"
	"net/http"
	"strconv"

	"bolao-copa/internal/pagination"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	page, limit := pagination.GetParams(c, 50)
	filters := MatchFilters{
		Stage:  c.Query("stage"),
		Status: c.Query("status"),
		Group:  c.Query("group"),
	}

	resp, err := h.service.ListMatches(c.Request.Context(), filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar jogos", "details": err.Error()})
		return
	}

	// 🛠️ INJETA OS GOLS GERADOS LOCALMENTE
	for i, m := range resp.Matches {
		resp.Matches[i].Goals = []Goal{} // Inicializa como array vazio contra null

		if m.Status == "IN_PLAY" || m.Status == "FINISHED" {
			// Se houver placar, gera os gols visualmente para o React
			homeScore := 0
			awayScore := 0
			if m.HomeScore != nil { homeScore = *m.HomeScore }
			if m.AwayScore != nil { awayScore = *m.AwayScore }

			// Simula placar do primeiro tempo (aproximadamente metade dos gols)
			homeHalfTime := homeScore / 2
			awayHalfTime := awayScore / 2

			resp.Matches[i].Goals = generateMockGoals(m.ID, m.HomeTeam.Name, m.AwayTeam.Name, homeHalfTime, awayHalfTime, homeScore, awayScore)
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	match, err := h.service.GetMatch(c.Request.Context(), id)
	if err == ErrMatchNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "jogo não encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar jogo"})
		return
	}

	// 🛠️ INJETA OS GOLS NO ENDPOINT DE ID TAMBÉM
	match.Goals = []Goal{}
	if match.Status == "IN_PLAY" || match.Status == "FINISHED" {
		homeScore := 0
		awayScore := 0
		if match.HomeScore != nil { homeScore = *match.HomeScore }
		if match.AwayScore != nil { awayScore = *match.AwayScore }

		// Simula placar do primeiro tempo (aproximadamente metade dos gols)
		homeHalfTime := homeScore / 2
		awayHalfTime := awayScore / 2

		match.Goals = generateMockGoals(match.ID, match.HomeTeam.Name, match.AwayTeam.Name, homeHalfTime, awayHalfTime, homeScore, awayScore)
	}

	c.JSON(http.StatusOK, match)
}

// 🪄 GERADOR DE GOLS AUTOMÁTICO BASEADO NO PLACAR DO BANCO
func generateMockGoals(matchID, homeTeam, awayTeam string, homeHalfTime, awayHalfTime, homeFullTime, awayFullTime int) []Goal {
	var mockGoals []Goal

	// Converte o ID para número para usar como semente simples de minutos
	idNum, _ := strconv.Atoi(matchID)
	if idNum == 0 {
		idNum = 10
	}

	// Calcula gols do 2º tempo (diferença entre fullTime e halfTime)
	homeSecondHalf := homeFullTime - homeHalfTime
	awaySecondHalf := awayFullTime - awayHalfTime

	// Gera gols do time da casa no 1º tempo (minutos 1-45)
	for i := 0; i < homeHalfTime; i++ {
		minuto := ((idNum + i*7) % 44) + 1 // Garante minuto entre 1-45
		mockGoals = append(mockGoals, Goal{
			PlayerName: fmt.Sprintf("Jogador %s %d", homeTeam, i+1),
			Minute:     minuto,
			Team:       "home",
		})
	}

	// Gera gols do time visitante no 1º tempo (minutos 1-45)
	for i := 0; i < awayHalfTime; i++ {
		minuto := ((idNum + i*11 + 5) % 44) + 1 // Garante minuto entre 1-45
		mockGoals = append(mockGoals, Goal{
			PlayerName: fmt.Sprintf("Jogador %s %d", awayTeam, i+1),
			Minute:     minuto,
			Team:       "away",
		})
	}

	// Gera gols do time da casa no 2º tempo (minutos 46-90)
	for i := 0; i < homeSecondHalf; i++ {
		minuto := ((idNum + i*9 + 20) % 45) + 46 // Garante minuto entre 46-90
		mockGoals = append(mockGoals, Goal{
			PlayerName: fmt.Sprintf("Jogador %s %d", homeTeam, homeHalfTime+i+1),
			Minute:     minuto,
			Team:       "home",
		})
	}

	// Gera gols do time visitante no 2º tempo (minutos 46-90)
	for i := 0; i < awaySecondHalf; i++ {
		minuto := ((idNum + i*13 + 25) % 45) + 46 // Garante minuto entre 46-90
		mockGoals = append(mockGoals, Goal{
			PlayerName: fmt.Sprintf("Jogador %s %d", awayTeam, awayHalfTime+i+1),
			Minute:     minuto,
			Team:       "away",
		})
	}

	return mockGoals
}