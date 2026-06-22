-- Adiciona colunas home_half_time e away_half_time na tabela matches
ALTER TABLE matches ADD COLUMN IF NOT EXISTS home_half_time INT;
ALTER TABLE matches ADD COLUMN IF NOT EXISTS away_half_time INT;