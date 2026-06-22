-- Adiciona coluna minute na tabela matches para armazenar a minutagem atual de jogos ao vivo
ALTER TABLE matches ADD COLUMN IF NOT EXISTS minute INT;
