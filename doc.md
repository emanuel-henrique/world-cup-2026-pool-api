# SRD — Bolão Copa do Mundo 2026
**Software Requirements Document**
Versão 1.2 | Autor: Emanuel | Data: Junho 2026

> **Changelog v1.2:** Estrutura de pastas atualizada para seguir convenções e padrões da comunidade Go (`cmd/`, `internal/`).
> **Changelog v1.1:** Adicionadas tabelas `players` e `groups_standing`; refatorada `special_predictions` para usar FK em vez de texto livre; adicionadas rotas de grupos, chaveamento e jogadores; atualizada estrutura de pastas.

---

## 1. Visão Geral

Aplicação web de bolão público para a Copa do Mundo 2026. Qualquer usuário pode se cadastrar, fazer palpites em jogos e competir em um ranking global. Os resultados são atualizados automaticamente via integração com API externa de dados esportivos, e a pontuação é calculada de forma automática ao fim de cada partida.

---

## 2. Objetivos

- Permitir que qualquer pessoa participe do bolão sem convite
- Cobrir toda a Copa: fase de grupos, mata-mata, artilheiro e campeão
- Calcular pontuação automaticamente com base nos resultados reais
- Manter um ranking global em tempo real
- Servir como projeto de portfólio demonstrando boas práticas em Go

---

## 3. Escopo

### 3.1 Dentro do escopo (v1.0)
- Cadastro e autenticação (email + senha)
- Palpites por jogo (placar exato)
- Palpite especial: artilheiro da Copa (lista de jogadores pré-definida)
- Palpite especial: campeão da Copa (lista de seleções)
- Visualização de grupos com tabela de classificação
- Visualização do chaveamento do mata-mata
- Fechamento automático de palpites no kickoff de cada jogo
- Cálculo automático de pontuação
- Ranking global de participantes
- Worker de sincronização com API externa (polling a cada 1 min)

### 3.2 Fora do escopo (v1.0)
- Bolões privados / grupos
- Chat ou comentários
- Notificações push
- App mobile nativo
- Pagamentos / premiação

---

## 4. Stakeholders

| Papel | Descrição |
|-------|-----------|
| Participante | Usuário cadastrado que faz palpites |
| Sistema | Worker Go que sincroniza dados e calcula pontuação |
| Administrador | Desenvolvedor com acesso direto ao banco |

---

## 5. Requisitos Funcionais

### 5.1 Autenticação

**RF-01** — O sistema deve permitir cadastro com nome, email e senha.
**RF-02** — O sistema deve autenticar usuários via email e senha, retornando um JWT.
**RF-03** — O sistema deve rejeitar cadastros com email já registrado.
**RF-04** — Senhas devem ser armazenadas com hash bcrypt.
**RF-05** — O JWT deve expirar em 7 dias.

### 5.2 Jogos

**RF-06** — O sistema deve exibir todos os jogos da Copa com data, horário, seleções e status.
**RF-07** — Os jogos devem ser sincronizados automaticamente via API externa (polling 1 min).
**RF-08** — Cada jogo deve ter um dos seguintes status: `scheduled`, `live`, `finished`.
**RF-09** — O resultado oficial deve ser persistido ao fim de cada jogo.

### 5.3 Grupos

**RF-10** — O sistema deve exibir os 12 grupos da Copa com suas seleções.
**RF-11** — Cada grupo deve exibir a tabela de classificação com pontos, saldo de gols, jogos disputados e aproveitamento.
**RF-12** — A tabela de classificação deve ser atualizada automaticamente a cada resultado.
**RF-13** — O sistema deve exibir os jogos de cada grupo separadamente.

### 5.4 Mata-mata

**RF-14** — O sistema deve exibir o chaveamento completo do mata-mata.
**RF-15** — O chaveamento deve cobrir: 32avos, oitavas, quartas, semifinais e final.
**RF-16** — Confrontos do mata-mata devem ser populados automaticamente pelo worker conforme a fase de grupos encerra.
**RF-17** — O chaveamento deve mostrar o vencedor de cada confronto quando disponível.

### 5.5 Palpites

**RF-18** — O usuário autenticado pode fazer palpite de placar exato em qualquer jogo com status `scheduled`.
**RF-19** — O sistema deve bloquear criação ou edição de palpites quando o jogo estiver `live` ou `finished`.
**RF-20** — O usuário pode editar o palpite até o kickoff do jogo.
**RF-21** — Cada usuário pode ter no máximo um palpite por jogo.
**RF-22** — O usuário pode escolher o artilheiro a partir de uma lista pré-definida de jogadores.
**RF-23** — O usuário pode escolher o campeão a partir da lista de seleções participantes.
**RF-24** — Palpites especiais (artilheiro e campeão) ficam abertos até o início do primeiro jogo da Copa.

### 5.6 Pontuação

**RF-25** — Ao fim de cada jogo, o sistema deve calcular a pontuação de todos os palpites daquele jogo automaticamente.
**RF-26** — A tabela de pontuação por jogo deve seguir as seguintes regras:

| Acerto | Pontos |
|--------|--------|
| Placar exato | 10 pts |
| Vencedor + saldo de gols correto | 7 pts |
| Vencedor correto | 5 pts |
| Empate previsto corretamente | 5 pts |
| Nenhum acerto | 0 pts |

**RF-27** — Palpite de campeão correto: **30 pontos**, concedidos após a final.
**RF-28** — Palpite de artilheiro correto: **20 pontos**, concedidos após a final.
**RF-29** — Pontuação total do usuário é a soma de todos os pontos acumulados.

### 5.7 Ranking

**RF-30** — O sistema deve exibir um ranking global com todos os participantes ordenados por pontuação total.
**RF-31** — Em caso de empate, o desempate é por maior número de placares exatos.
**RF-32** — O ranking deve ser atualizado a cada vez que uma pontuação é calculada.

---

## 6. Requisitos Não Funcionais

**RNF-01** — A API Go deve responder em menos de 300ms para endpoints de leitura.
**RNF-02** — O worker de polling não deve fazer mais de 60 requisições/hora à API externa (respeitar limite do tier gratuito).
**RNF-03** — Todos os endpoints protegidos devem exigir JWT válido.
**RNF-04** — O sistema deve funcionar com pelo menos 100 usuários simultâneos sem degradação.
**RNF-05** — O banco de dados deve ser PostgreSQL (Neon).
**RNF-06** — O backend deve ser deployado no Render (free tier).
**RNF-07** — O frontend deve ser deployado na Vercel.
**RNF-08** — Logs de erro do worker devem ser registrados com contexto suficiente para debug.

---

## 7. Regras de Negócio

**RN-01** — Um palpite só pode ser registrado se o jogo ainda não começou (status `scheduled`).
**RN-02** — O horário de corte para palpites é o `kickoff_at` do jogo, validado no servidor (não no cliente).
**RN-03** — Palpites especiais (artilheiro e campeão) são bloqueados no início do primeiro jogo da Copa.
**RN-04** — O cálculo de pontuação é disparado pelo worker quando detecta que o status mudou para `finished`.
**RN-05** — Pontuação já calculada não pode ser recalculada (idempotência via flag `scored` no palpite).
**RN-06** — O resultado do jogo vem da API externa e é considerado fonte da verdade.
**RN-07** — Empate na fase de grupos conta como empate. No mata-mata, o resultado final (incluindo prorrogação e pênaltis) é o que vale.
**RN-08** — A classificação de cada grupo é calculada dinamicamente via query — não é armazenada como tabela separada.
**RN-09** — O chaveamento do mata-mata só é populado após o encerramento completo da fase de grupos.

---

## 8. Arquitetura

### 8.1 Visão Geral

```
[React/TypeScript Frontend]
        |
        | HTTP/JSON
        v
[Go HTTP Server — Render]
        |
   +----|----+
   |         |
[PostgreSQL]  [API-Football Worker]
  (Neon)      (time.Ticker — 1 min)
```

### 8.2 Backend (Go)

- Estrutura flat por entidade (sem DDD, sem camadas forçadas)
- Funções planas, sem struct methods
- Interface de repositório por entidade (permite mock em testes)
- `database/sql` com `lib/pq`
- JWT via `golang-jwt/jwt`
- Sem frameworks — `net/http` puro com roteador simples
- Variáveis de ambiente via `os.Getenv`

### 8.3 Worker

- Roda em goroutine separada na mesma instância Go
- `time.Ticker` com intervalo de 1 minuto
- Busca jogos com status `scheduled` ou `live` na API externa
- Atualiza status e placar no banco
- Quando detecta `finished`, dispara cálculo de pontuação para aquele jogo
- Popula confrontos do mata-mata quando fase de grupos encerra
- Usa flag `scored` nos palpites para garantir idempotência

### 8.4 Frontend (React/TypeScript)

- Vite + React + TypeScript
- Tailwind CSS + shadcn/ui
- Polling leve no cliente (a cada 60s) para atualizar ranking e jogos
- Sem WebSocket (simplicidade > complexidade)

---

## 9. Schema do Banco de Dados

```sql
-- Usuários
CREATE TABLE users (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT NOT NULL,
  email      TEXT UNIQUE NOT NULL,
  password   TEXT NOT NULL, -- bcrypt hash
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Seleções
CREATE TABLE teams (
  id   TEXT PRIMARY KEY, -- ex: "BRA", "ARG"
  name TEXT NOT NULL,
  flag TEXT -- URL da bandeira
);

-- Jogadores (para lista de artilheiros)
CREATE TABLE players (
  id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name    TEXT NOT NULL,
  team_id TEXT NOT NULL REFERENCES teams(id)
);

-- Jogos
CREATE TABLE matches (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  external_id  TEXT UNIQUE NOT NULL,       -- ID da API externa
  home_team_id TEXT REFERENCES teams(id),  -- NULL no mata-mata até classificados definidos
  away_team_id TEXT REFERENCES teams(id),
  home_score   INT,
  away_score   INT,
  stage        TEXT NOT NULL,              -- "group" | "round_of_32" | "round_of_16" | "quarter" | "semi" | "final"
  group_name   TEXT,                       -- "A"..."L" (NULL no mata-mata)
  kickoff_at   TIMESTAMPTZ NOT NULL,
  status       TEXT NOT NULL DEFAULT 'scheduled', -- scheduled | live | finished
  created_at   TIMESTAMPTZ DEFAULT NOW()
);

-- Palpites por jogo
CREATE TABLE predictions (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id),
  match_id   UUID NOT NULL REFERENCES matches(id),
  home_score INT NOT NULL,
  away_score INT NOT NULL,
  points     INT NOT NULL DEFAULT 0,
  scored     BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(user_id, match_id)
);

-- Palpites especiais
CREATE TABLE special_predictions (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        UUID NOT NULL REFERENCES users(id) UNIQUE,
  champion_id    TEXT REFERENCES teams(id),   -- seleção campeã
  top_scorer_id  UUID REFERENCES players(id), -- artilheiro (lista pré-definida)
  points         INT NOT NULL DEFAULT 0,
  scored         BOOLEAN NOT NULL DEFAULT FALSE,
  created_at     TIMESTAMPTZ DEFAULT NOW(),
  updated_at     TIMESTAMPTZ DEFAULT NOW()
);

-- Classificação de grupos: calculada via VIEW, não tabela
-- Evita duplicidade de dados e inconsistência
CREATE VIEW group_standings AS
SELECT
  m.group_name,
  t.id   AS team_id,
  t.name AS team_name,
  t.flag,
  COUNT(*)                                                        AS played,
  SUM(CASE
    WHEN (m.home_team_id = t.id AND m.home_score > m.away_score)
      OR (m.away_team_id = t.id AND m.away_score > m.home_score) THEN 1 ELSE 0
  END)                                                            AS wins,
  SUM(CASE WHEN m.home_score = m.away_score THEN 1 ELSE 0 END)   AS draws,
  SUM(CASE
    WHEN (m.home_team_id = t.id AND m.home_score < m.away_score)
      OR (m.away_team_id = t.id AND m.away_score < m.home_score) THEN 1 ELSE 0
  END)                                                            AS losses,
  SUM(CASE
    WHEN m.home_team_id = t.id THEN m.home_score
    WHEN m.away_team_id = t.id THEN m.away_score ELSE 0
  END)                                                            AS goals_for,
  SUM(CASE
    WHEN m.home_team_id = t.id THEN m.away_score
    WHEN m.away_team_id = t.id THEN m.home_score ELSE 0
  END)                                                            AS goals_against,
  SUM(CASE
    WHEN (m.home_team_id = t.id AND m.home_score > m.away_score)
      OR (m.away_team_id = t.id AND m.away_score > m.home_score) THEN 3
    WHEN m.home_score = m.away_score THEN 1
    ELSE 0
  END)                                                            AS points
FROM matches m
JOIN teams t ON t.id IN (m.home_team_id, m.away_team_id)
WHERE m.stage = 'group'
  AND m.status = 'finished'
  AND m.group_name IS NOT NULL
GROUP BY m.group_name, t.id, t.name, t.flag
ORDER BY m.group_name, points DESC, (goals_for - goals_against) DESC, goals_for DESC;
```

---

## 10. Endpoints da API

### Autenticação
| Método | Rota | Descrição | Auth |
|--------|------|-----------|------|
| POST | `/auth/register` | Cadastro | ❌ |
| POST | `/auth/login` | Login, retorna JWT | ❌ |

### Jogos
| Método | Rota | Descrição | Auth |
|--------|------|-----------|------|
| GET | `/matches` | Lista todos os jogos | ❌ |
| GET | `/matches/:id` | Detalhe de um jogo | ❌ |

### Grupos
| Método | Rota | Descrição | Auth |
|--------|------|-----------|------|
| GET | `/groups` | Lista os 12 grupos com classificação | ❌ |
| GET | `/groups/:name` | Detalhe do grupo: classificação + jogos | ❌ |

### Mata-mata
| Método | Rota | Descrição | Auth |
|--------|------|-----------|------|
| GET | `/bracket` | Chaveamento completo do mata-mata | ❌ |

### Jogadores
| Método | Rota | Descrição | Auth |
|--------|------|-----------|------|
| GET | `/players` | Lista jogadores disponíveis para palpite de artilheiro | ❌ |

### Palpites
| Método | Rota | Descrição | Auth |
|--------|------|-----------|------|
| GET | `/predictions` | Meus palpites de jogos | ✅ |
| POST | `/predictions` | Criar/atualizar palpite de jogo | ✅ |
| GET | `/predictions/special` | Meu palpite especial | ✅ |
| POST | `/predictions/special` | Criar/atualizar palpite especial | ✅ |

### Ranking
| Método | Rota | Descrição | Auth |
|--------|------|-----------|------|
| GET | `/ranking` | Ranking global paginado | ❌ |

---

## 11. Estrutura de Pastas (Go)

Segue as convenções da comunidade Go: `cmd/` para entrypoints e `internal/` para todo código da aplicação. O compilador Go impede que pacotes dentro de `internal/` sejam importados por módulos externos — isolamento garantido em tempo de compilação.

```
bolao-copa/
├── cmd/
│   └── api/
│       └── main.go              -- entrypoint: inicializa DB, rotas e worker
├── internal/
│   ├── auth/
│   │   ├── handler.go           -- HTTP: register, login
│   │   ├── middleware.go        -- JWT middleware
│   │   ├── repository.go        -- interface + implementação PostgreSQL
│   │   └── service.go           -- lógica: depende da interface, não do banco
│   ├── matches/
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── groups/
│   │   ├── handler.go
│   │   └── service.go           -- consulta a VIEW group_standings
│   ├── bracket/
│   │   ├── handler.go
│   │   └── service.go           -- monta o chaveamento do mata-mata
│   ├── players/
│   │   ├── handler.go
│   │   └── repository.go
│   ├── predictions/
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── ranking/
│   │   ├── handler.go
│   │   └── service.go
│   ├── scoring/
│   │   └── scoring.go           -- lógica pura de pontuação (testável sem banco)
│   ├── worker/
│   │   └── sync.go              -- polling + trigger de pontuação + popula mata-mata
│   └── db/
│       └── db.go                -- conexão com PostgreSQL
├── db/
│   └── migrations/
│       └── 001_init.sql
├── .env
├── go.mod
└── go.sum
```

---

## 12. Fluxo do Worker

```
A cada 1 minuto:
  1. Buscar na API-Football todos os jogos da Copa com status != finished
  2. Para cada jogo retornado:
     a. Comparar com o registro no banco
     b. Se status mudou → atualizar banco
     c. Se status virou "finished" e jogo ainda não foi pontuado:
        → chamar scoring.CalculateMatch(matchID)
        → marcar todos os predictions daquele jogo como scored=true
  3. Verificar se todos os jogos da fase de grupos encerraram:
     → Se sim, popular confrontos do mata-mata (round_of_32) com os classificados
  4. Registrar log com resumo da execução
```

---

## 13. Critérios de Aceite (MVP)

- [ ] Usuário consegue se cadastrar e logar
- [ ] Usuário vê lista de todos os jogos com status e horário
- [ ] Usuário vê tabela de classificação de cada grupo
- [ ] Usuário vê chaveamento do mata-mata atualizado
- [ ] Usuário escolhe artilheiro a partir de lista de jogadores
- [ ] Usuário escolhe campeão a partir de lista de seleções
- [ ] Usuário consegue fazer palpite antes do kickoff
- [ ] Sistema bloqueia palpite após kickoff
- [ ] Worker sincroniza resultado após fim do jogo
- [ ] Pontuação é calculada automaticamente
- [ ] Ranking exibe todos os usuários ordenados por pontos
- [ ] Palpite especial (artilheiro + campeão) funciona e pontua ao fim da Copa

---

## 14. Tecnologias

| Componente | Tecnologia |
|------------|------------|
| Backend | Go (net/http) |
| Banco | PostgreSQL — Neon |
| Auth | JWT (golang-jwt/jwt) |
| Hash | bcrypt |
| Frontend | React + TypeScript + Vite |
| Estilo | Tailwind CSS + shadcn/ui |
| Deploy Backend | Render |
| Deploy Frontend | Vercel |
| Dados esportivos | API-Football (RapidAPI, tier gratuito) |

---

*Documento gerado como base para desenvolvimento. Sujeito a revisão conforme evolução do projeto.*
