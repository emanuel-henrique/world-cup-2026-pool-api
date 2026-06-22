// internal/worker/sync.go
package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"bolao-copa/internal/scoring"
)

type Worker struct {
	db     *sql.DB
	client *http.Client
}

func New(db *sql.DB) *Worker {
	return &Worker{
		db:     db,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *Worker) Start(ctx context.Context) {
	log.Println("[worker] iniciado — aguardando alinhamento com o próximo minuto exato")

	// Aguarda até o próximo minuto exato (segundo 00)
	nextMinute := time.Now().Truncate(time.Minute).Add(time.Minute)
	sleepDuration := time.Until(nextMinute)
	log.Printf("[worker] aguardando %v até o próximo minuto exato", sleepDuration)
	
	select {
	case <-time.After(sleepDuration):
		log.Println("[worker] alinhado com o minuto exato — iniciando loop de sincronização")
	case <-ctx.Done():
		log.Println("[worker] encerrado antes do alinhamento")
		return
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Roda imediatamente após o alinhamento
	w.run(ctx)

	for {
		select {
		case <-ticker.C:
			w.run(ctx)
		case <-ctx.Done():
			log.Println("[worker] encerrado")
			return
		}
	}
}

func (w *Worker) run(ctx context.Context) {
	log.Println("[worker] iniciando sincronização...")

	fixtures, err := w.fetchFixtures()
	if err != nil {
		log.Printf("[worker] erro ao buscar jogos da API: %v", err)
		return
	}

	for _, f := range fixtures {
		prev, err := w.getMatchStatus(ctx, f.ExternalID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("[worker] erro ao buscar status do jogo %s: %v", f.ExternalID, err)
			continue
		}

		// Se o jogo não existe, cria
		if prev == "" {
			err = w.createMatch(ctx, f)
			if err != nil {
				log.Printf("[worker] erro ao criar jogo %s: %v", f.ExternalID, err)
				continue
			}
		} else {
			// Atualiza status e placar no banco
			err = w.updateMatch(ctx, f)
			if err != nil {
				log.Printf("[worker] erro ao atualizar jogo %s: %v", f.ExternalID, err)
				continue
			}
		}

		// Se acabou de terminar → calcula pontuação
		if prev != "finished" && f.Status == "finished" {
			matchID, err := w.getMatchID(ctx, f.ExternalID)
			if err != nil {
				log.Printf("[worker] erro ao buscar id do jogo %s: %v", f.ExternalID, err)
				continue
			}

			if err := scoring.CalculateMatch(ctx, w.db, matchID); err != nil {
				log.Printf("[worker] erro ao calcular pontuação do jogo %s: %v", matchID, err)
				continue
			}

			log.Printf("[worker] pontuação calculada para o jogo %s", matchID)
		}
	}

	// Verifica se fase de grupos encerrou
	if err := w.populateBracketIfReady(ctx); err != nil {
		log.Printf("[worker] erro ao popular mata-mata: %v", err)
	}

	log.Println("[worker] sincronização concluída")
}

// --- API-Football ---

type Goal struct {
	PlayerName string `json:"playerName"`
	Minute     int    `json:"minute"`
	Team       string `json:"team"` // "home" or "away"
}

type fixture struct {
	ExternalID string
	HomeTeamID string
	AwayTeamID string
	HomeScore    *int
	AwayScore    *int
	HomeHalfTime *int
	AwayHalfTime *int
	Status       string
	KickoffAt  time.Time
	Stage      string
	GroupName  *string
	Goals      []Goal
	Minute     *int // current minute for live matches
}

type apiResponse struct {
	Matches []struct {
		ID      int    `json:"id"`
		UtcDate string `json:"utcDate"`
		Status  string `json:"status"`
		Stage   string `json:"stage"`
		Group   string `json:"group"`
		HomeTeam struct {
			ID int `json:"id"`
		} `json:"homeTeam"`
		AwayTeam struct {
			ID int `json:"id"`
		} `json:"awayTeam"`
		Score struct {
			FullTime struct {
				Home *int `json:"home"`
				Away *int `json:"away"`
			} `json:"fullTime"`
			HalfTime struct {
				Home *int `json:"home"`
				Away *int `json:"away"`
			} `json:"halfTime"`
		} `json:"score"`
		Goals []struct {
			PlayerName string `json:"playerName"`
			Minute     int    `json:"minute"`
			Team       string `json:"team"` // "home" or "away"
		} `json:"goals,omitempty"`
		Events []struct {
			PlayerName string `json:"playerName"`
			Minute     int    `json:"minute"`
			Team       string `json:"team"` // "home" or "away"
			Type       string `json:"type"` // "goal", "card", etc.
		} `json:"events,omitempty"`
		Minute *int `json:"minute"` // current minute for live matches
	} `json:"matches"`
}

func (w *Worker) fetchFixtures() ([]fixture, error) {
	apiKey := os.Getenv("API_FOOTBALL_KEY")
	url := "https://api.football-data.org/v4/competitions/WC/matches"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}
	req.Header.Set("X-Auth-Token", apiKey)
	req.Header.Set("X-Unfold-Goals", "true")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na API: status %d", resp.StatusCode)
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	var fixtures []fixture
	for _, m := range data.Matches {
		status := MapStatus(m.Status)
		kickoff, _ := time.Parse(time.RFC3339, m.UtcDate)

		stage, groupName := MapStageAndGroup(m.Stage, m.Group)

		// Process goals - tenta tanto de Goals quanto de Events
		var goals []Goal
		if m.Goals != nil && len(m.Goals) > 0 {
			for _, g := range m.Goals {
				if g.PlayerName != "" && g.Team != "" {
					goals = append(goals, Goal{
						PlayerName: g.PlayerName,
						Minute:     g.Minute,
						Team:       g.Team,
					})
				}
			}
		}
		if m.Events != nil && len(m.Events) > 0 {
			for _, e := range m.Events {
				if e.Type == "goal" && e.PlayerName != "" && e.Team != "" {
					goals = append(goals, Goal{
						PlayerName: e.PlayerName,
						Minute:     e.Minute,
						Team:       e.Team,
					})
				}
			}
		}

		f := fixture{
			ExternalID:   fmt.Sprintf("%d", m.ID),
			HomeTeamID:   fmt.Sprintf("%d", m.HomeTeam.ID),
			AwayTeamID:   fmt.Sprintf("%d", m.AwayTeam.ID),
			HomeScore:    m.Score.FullTime.Home,
			AwayScore:    m.Score.FullTime.Away,
			HomeHalfTime: m.Score.HalfTime.Home, // Adicionado
			AwayHalfTime: m.Score.HalfTime.Away, // Adicionado
			Status:       status,
			KickoffAt:    kickoff,
			Stage:        stage,
			GroupName:    groupName,
			Goals:        goals,
			Minute:       m.Minute,
		}
		fixtures = append(fixtures, f)
	}

	return fixtures, nil
}


func MapStatus(status string) string {
	switch status {
	case "SCHEDULED", "TIMED":
		return "scheduled"
	case "IN_PLAY", "PAUSED":
		return "live"
	case "FINISHED", "AWARDED":
		return "finished"
	default:
		return "scheduled"
	}
}

func MapStageAndGroup(stage, group string) (string, *string) {
	var mappedStage string
	switch stage {
	case "GROUP_STAGE":
		mappedStage = "group"
	case "LAST_32":
		mappedStage = "round_of_32"
	case "LAST_16":
		mappedStage = "round_of_16"
	case "QUARTER_FINALS":
		mappedStage = "quarter"
	case "SEMI_FINALS":
		mappedStage = "semi"
	case "FINAL":
		mappedStage = "final"
	case "THIRD_PLACE":
		mappedStage = "third_place"
	default:
		mappedStage = "group"
	}

	var groupName *string
	if len(group) > 6 && group[:6] == "GROUP_" {
		name := group[6:]
		groupName = &name
	}

	return mappedStage, groupName
}

// --- Helpers de banco ---

func (w *Worker) getMatchStatus(ctx context.Context, externalID string) (string, error) {
	var status string
	err := w.db.QueryRowContext(ctx,
		`SELECT status FROM matches WHERE external_id = $1`, externalID,
	).Scan(&status)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return status, err
}

func (w *Worker) getMatchID(ctx context.Context, externalID string) (string, error) {
	var id string
	err := w.db.QueryRowContext(ctx,
		`SELECT id FROM matches WHERE external_id = $1`, externalID,
	).Scan(&id)
	return id, err
}

func (w *Worker) createMatch(ctx context.Context, f fixture) error {
    // Primeiro, buscar os IDs dos times no banco
    var homeTeamID, awayTeamID string
    err := w.db.QueryRowContext(ctx, `SELECT id FROM teams WHERE external_id = $1`, f.HomeTeamID).Scan(&homeTeamID)
    if err != nil {
        return fmt.Errorf("erro ao buscar time home: %w", err)
    }
    err = w.db.QueryRowContext(ctx, `SELECT id FROM teams WHERE external_id = $1`, f.AwayTeamID).Scan(&awayTeamID)
    if err != nil {
        return fmt.Errorf("erro ao buscar time away: %w", err)
    }

    // Criar o jogo
    var matchID string
    err = w.db.QueryRowContext(ctx, `
        INSERT INTO matches (external_id, home_team_id, away_team_id, home_score, away_score, home_half_time, away_half_time, stage, group_name, kickoff_at, status, minute)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id
    `, f.ExternalID, homeTeamID, awayTeamID, f.HomeScore, f.AwayScore, f.HomeHalfTime, f.AwayHalfTime, f.Stage, f.GroupName, f.KickoffAt, f.Status, f.Minute).Scan(&matchID)
    if err != nil {
        return fmt.Errorf("erro ao criar jogo: %w", err)
    }

    // Inserir gols
    for _, g := range f.Goals {
        _, err = w.db.ExecContext(ctx, `
            INSERT INTO goals (match_id, player_name, minute, team)
            VALUES ($1, $2, $3, $4)
        `, matchID, g.PlayerName, g.Minute, g.Team)
        if err != nil {
            log.Printf("[worker] erro ao inserir gol: %v", err)
            continue
        }
    }

    return nil
}

func (w *Worker) updateMatch(ctx context.Context, f fixture) error {
    _, err := w.db.ExecContext(ctx, `
        UPDATE matches
        SET status         = $1,
            home_score     = $2,
            away_score     = $3,
            home_half_time = $4,
            away_half_time = $5,
            minute         = $6
        WHERE external_id = $7
    `, f.Status, f.HomeScore, f.AwayScore, f.HomeHalfTime, f.AwayHalfTime, f.Minute, f.ExternalID)
    if err != nil {
        return err
    }

    // Salvar gols no banco
    matchID, err := w.getMatchID(ctx, f.ExternalID)
    if err != nil {
        return fmt.Errorf("erro ao buscar id do jogo: %w", err)
    }

    // Deletar gols antigos deste match
    _, err = w.db.ExecContext(ctx, `DELETE FROM goals WHERE match_id = $1`, matchID)
    if err != nil {
        return fmt.Errorf("erro ao deletar gols antigos: %w", err)
    }

    // Inserir novos gols
    for _, g := range f.Goals {
        _, err = w.db.ExecContext(ctx, `
            INSERT INTO goals (match_id, player_name, minute, team)
            VALUES ($1, $2, $3, $4)
        `, matchID, g.PlayerName, g.Minute, g.Team)
        if err != nil {
            log.Printf("[worker] erro ao inserir gol: %v", err)
            continue
        }
    }

    return nil
}

func (w *Worker) populateBracketIfReady(ctx context.Context) error {
	// Verifica se todos os jogos da fase de grupos encerraram
	var pending int
	err := w.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM matches
		WHERE stage = 'group' AND status != 'finished'
	`).Scan(&pending)
	if err != nil {
		return fmt.Errorf("erro ao verificar fase de grupos: %w", err)
	}

	if pending > 0 {
		return nil // fase de grupos ainda não encerrou
	}

	// Verifica se mata-mata já foi populado
	var bracket int
	err = w.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM matches WHERE stage = 'round_of_32'
	`).Scan(&bracket)
	if err != nil {
		return fmt.Errorf("erro ao verificar mata-mata: %w", err)
	}

	if bracket > 0 {
		return nil // mata-mata já populado
	}

	log.Println("[worker] fase de grupos encerrada — populando mata-mata...")

	// Por ora loga — a população real depende das regras de classificação da FIFA
	// que serão implementadas junto com o seed de dados
	log.Println("[worker] TODO: popular confrontos do round_of_32")

	return nil
}
