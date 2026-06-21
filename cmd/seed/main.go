// cmd/seed/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"bolao-copa/internal/database"
	"bolao-copa/internal/worker"
)

type TeamsResponse struct {
	Teams []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Crest string `json:"crest"`
		Squad []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"squad"`
	} `json:"teams"`
}

type MatchesResponse struct {
	Matches []struct {
		ID       int    `json:"id"`
		UtcDate  string `json:"utcDate"`
		Status   string `json:"status"`
		Stage    string `json:"stage"`
		Group    string `json:"group"`
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
		} `json:"score"`
	} `json:"matches"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("aviso: .env não encontrado, usando variáveis de ambiente")
	}

	apiKey := os.Getenv("API_FOOTBALL_KEY")
	if apiKey == "" {
		log.Fatal("API_FOOTBALL_KEY não definida")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("falha ao conectar no banco: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	client := &http.Client{Timeout: 30 * time.Second}

	log.Println("[seeder] Baixando seleções (teams)...")
	teamsReq, err := http.NewRequest("GET", "https://api.football-data.org/v4/competitions/WC/teams", nil)
	if err != nil {
		log.Fatalf("erro request teams: %v", err)
	}
	teamsReq.Header.Set("X-Auth-Token", apiKey)

	teamsResp, err := client.Do(teamsReq)
	if err != nil {
		log.Fatalf("erro chamada teams: %v", err)
	}
	defer teamsResp.Body.Close()

	var teamsData TeamsResponse
	if err := json.NewDecoder(teamsResp.Body).Decode(&teamsData); err != nil {
		log.Fatalf("erro decode teams: %v", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	for _, t := range teamsData.Teams {
		idStr := fmt.Sprintf("%d", t.ID)
		_, err := tx.ExecContext(ctx, `
			INSERT INTO teams (id, name, flag)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, flag = EXCLUDED.flag
		`, idStr, t.Name, t.Crest)
		if err != nil {
			log.Printf("erro ao inserir time %s: %v", t.Name, err)
		}

		for _, p := range t.Squad {
			pIDStr := fmt.Sprintf("%d", p.ID)
			_, err := tx.ExecContext(ctx, `
				INSERT INTO players (id, name, team_id)
				VALUES ($1, $2, $3)
				ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, team_id = EXCLUDED.team_id
			`, pIDStr, p.Name, idStr)
			if err != nil {
				log.Printf("erro ao inserir jogador %s: %v", p.Name, err)
			}
		}
	}
	log.Printf("[seeder] %d seleções sincronizadas.", len(teamsData.Teams))

	log.Println("[seeder] Baixando jogos (matches)...")
	matchesReq, err := http.NewRequest("GET", "https://api.football-data.org/v4/competitions/WC/matches", nil)
	if err != nil {
		log.Fatalf("erro request matches: %v", err)
	}
	matchesReq.Header.Set("X-Auth-Token", apiKey)

	matchesResp, err := client.Do(matchesReq)
	if err != nil {
		log.Fatalf("erro chamada matches: %v", err)
	}
	defer matchesResp.Body.Close()

	var matchesData MatchesResponse
	if err := json.NewDecoder(matchesResp.Body).Decode(&matchesData); err != nil {
		log.Fatalf("erro decode matches: %v", err)
	}

	for _, m := range matchesData.Matches {
		stage, groupName := worker.MapStageAndGroup(m.Stage, m.Group)
		status := worker.MapStatus(m.Status)
		kickoff, _ := time.Parse(time.RFC3339, m.UtcDate)
		extID := fmt.Sprintf("%d", m.ID)

		var homeID, awayID *string
		if m.HomeTeam.ID != 0 {
			h := fmt.Sprintf("%d", m.HomeTeam.ID)
			homeID = &h
		}
		if m.AwayTeam.ID != 0 {
			a := fmt.Sprintf("%d", m.AwayTeam.ID)
			awayID = &a
		}

		_, err := tx.ExecContext(ctx, `
			INSERT INTO matches (external_id, home_team_id, away_team_id, home_score, away_score, stage, group_name, kickoff_at, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (external_id) DO UPDATE SET
				home_team_id = EXCLUDED.home_team_id,
				away_team_id = EXCLUDED.away_team_id,
				home_score = EXCLUDED.home_score,
				away_score = EXCLUDED.away_score,
				stage = EXCLUDED.stage,
				group_name = EXCLUDED.group_name,
				kickoff_at = EXCLUDED.kickoff_at,
				status = EXCLUDED.status
		`, extID, homeID, awayID, m.Score.FullTime.Home, m.Score.FullTime.Away, stage, groupName, kickoff, status)
		
		if err != nil {
			log.Printf("erro ao inserir match %s: %v", extID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("erro ao commitar transação: %v", err)
	}

	log.Printf("[seeder] %d jogos sincronizados.", len(matchesData.Matches))
	log.Println("[seeder] Concluído com sucesso!")
}
