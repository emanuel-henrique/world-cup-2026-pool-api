CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Usuários
CREATE TABLE IF NOT EXISTS users (
  id         TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
  name       TEXT NOT NULL,
  email      TEXT UNIQUE NOT NULL,
  password   TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seleções
CREATE TABLE IF NOT EXISTS teams (
  id   TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  flag TEXT
);

-- Jogadores
CREATE TABLE IF NOT EXISTS players (
  id      TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
  name    TEXT NOT NULL,
  team_id TEXT NOT NULL REFERENCES teams(id)
);

-- Jogos
CREATE TABLE IF NOT EXISTS matches (
  id           TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
  external_id  TEXT UNIQUE NOT NULL,
  home_team_id TEXT REFERENCES teams(id),
  away_team_id TEXT REFERENCES teams(id),
  home_score   INT,
  away_score   INT,
  stage        TEXT NOT NULL,
  group_name   TEXT,
  kickoff_at   TIMESTAMPTZ NOT NULL,
  status       TEXT NOT NULL DEFAULT 'scheduled',
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Palpites por jogo
CREATE TABLE IF NOT EXISTS predictions (
  id         TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
  user_id    TEXT NOT NULL REFERENCES users(id),
  match_id   TEXT NOT NULL REFERENCES matches(id),
  home_score INT NOT NULL,
  away_score INT NOT NULL,
  points     INT NOT NULL DEFAULT 0,
  scored     BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, match_id)
);

-- Palpites especiais
CREATE TABLE IF NOT EXISTS special_predictions (
  id            TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
  user_id       TEXT NOT NULL REFERENCES users(id) UNIQUE,
  champion_id   TEXT REFERENCES teams(id),
  top_scorer_id TEXT REFERENCES players(id),
  points        INT NOT NULL DEFAULT 0,
  scored        BOOLEAN NOT NULL DEFAULT FALSE,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- View de classificação de grupos
CREATE OR REPLACE VIEW group_standings AS
SELECT
  m.group_name,
  t.id   AS team_id,
  t.name AS team_name,
  t.flag,
  COUNT(*)                                                          AS played,
  SUM(CASE
    WHEN (m.home_team_id = t.id AND m.home_score > m.away_score)
      OR (m.away_team_id = t.id AND m.away_score > m.home_score)
    THEN 1 ELSE 0
  END)                                                              AS wins,
  SUM(CASE WHEN m.home_score = m.away_score THEN 1 ELSE 0 END)     AS draws,
  SUM(CASE
    WHEN (m.home_team_id = t.id AND m.home_score < m.away_score)
      OR (m.away_team_id = t.id AND m.away_score < m.home_score)
    THEN 1 ELSE 0
  END)                                                              AS losses,
  SUM(CASE
    WHEN m.home_team_id = t.id THEN m.home_score
    WHEN m.away_team_id = t.id THEN m.away_score
    ELSE 0
  END)                                                              AS goals_for,
  SUM(CASE
    WHEN m.home_team_id = t.id THEN m.away_score
    WHEN m.away_team_id = t.id THEN m.home_score
    ELSE 0
  END)                                                              AS goals_against,
  SUM(CASE
    WHEN (m.home_team_id = t.id AND m.home_score > m.away_score)
      OR (m.away_team_id = t.id AND m.away_score > m.home_score)
    THEN 3
    WHEN m.home_score = m.away_score THEN 1
    ELSE 0
  END)                                                              AS points
FROM matches m
JOIN teams t ON t.id IN (m.home_team_id, m.away_team_id)
WHERE m.stage = 'group'
  AND m.status = 'finished'
  AND m.group_name IS NOT NULL
GROUP BY m.group_name, t.id, t.name, t.flag
ORDER BY m.group_name, points DESC, (goals_for - goals_against) DESC, goals_for DESC;