-- Cria tabela goals para armazenar os gols de cada partida
CREATE TABLE IF NOT EXISTS goals (
  id          TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
  match_id    TEXT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
  player_name TEXT NOT NULL,
  minute      INT NOT NULL,
  team        TEXT NOT NULL, -- "home" or "away"
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
