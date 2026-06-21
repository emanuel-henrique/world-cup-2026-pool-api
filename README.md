# ⚽ Bolão Copa 2026

API REST para um bolão público da Copa do Mundo 2026. Qualquer pessoa pode se cadastrar, fazer palpites nos jogos e competir em um ranking global. Os resultados são sincronizados automaticamente com uma API externa e a pontuação é calculada sem intervenção manual.

> Projeto de portfólio desenvolvido com foco em boas práticas de Go, princípios SOLID e arquitetura limpa.

---

## 🚀 Stack

| Camada | Tecnologia |
|--------|------------|
| Backend | Go + Gin |
| Banco de dados | PostgreSQL (Neon) |
| Autenticação | JWT (golang-jwt/jwt) + bcrypt |
| Dados esportivos | API-Football (RapidAPI) |
| Deploy | Render |
| Frontend | React + TypeScript + Vite (repositório separado) |

---

## ✨ Funcionalidades

- Cadastro e autenticação via email e senha
- Palpites de placar exato por jogo
- Palpite especial: artilheiro e campeão da Copa
- Tabela de classificação de cada grupo (calculada via SQL VIEW)
- Chaveamento do mata-mata atualizado automaticamente
- Fechamento automático de palpites no kickoff de cada jogo
- Pontuação calculada automaticamente ao fim de cada partida
- Ranking global com critério de desempate por placares exatos
- Worker de sincronização com API externa (polling a cada 1 min)

---

## 🏆 Tabela de Pontuação

| Acerto | Pontos |
|--------|--------|
| Placar exato | 10 pts |
| Vencedor + saldo de gols correto | 7 pts |
| Vencedor correto | 5 pts |
| Empate previsto corretamente | 5 pts |
| Campeão correto | 30 pts |
| Artilheiro correto | 20 pts |
| Nenhum acerto | 0 pts |

---

## 📁 Estrutura do Projeto

```
bolao-copa/
├── cmd/
│   └── api/
│       └── main.go              # entrypoint
├── internal/
│   ├── auth/                    # autenticação, JWT, middleware
│   ├── matches/                 # jogos e sincronização
│   ├── groups/                  # classificação de grupos
│   ├── bracket/                 # chaveamento do mata-mata
│   ├── players/                 # lista de jogadores
│   ├── predictions/             # palpites por jogo e especiais
│   ├── ranking/                 # ranking global
│   ├── scoring/                 # lógica pura de pontuação
│   ├── worker/                  # polling da API externa
│   ├── router/                  # registro de rotas Gin
│   └── db/                      # conexão com PostgreSQL
├── db/
│   └── migrations/
│       └── 001_init.sql
├── .env.example
├── go.mod
└── go.sum
```

Cada módulo segue a mesma estrutura interna:

```
internal/<módulo>/
├── model.go        # structs de domínio
├── dto.go          # structs de entrada/saída da API
├── repository.go   # interface + implementação PostgreSQL
├── service.go      # lógica de negócio
└── handler.go      # handlers HTTP
```

---

## 🔌 Endpoints

### Autenticação
| Método | Rota | Auth |
|--------|------|------|
| POST | `/auth/register` | ❌ |
| POST | `/auth/login` | ❌ |

### Jogos
| Método | Rota | Auth |
|--------|------|------|
| GET | `/matches` | ❌ |
| GET | `/matches/:id` | ❌ |

### Grupos
| Método | Rota | Auth |
|--------|------|------|
| GET | `/groups` | ❌ |
| GET | `/groups/:name` | ❌ |

### Mata-mata
| Método | Rota | Auth |
|--------|------|------|
| GET | `/bracket` | ❌ |

### Jogadores
| Método | Rota | Auth |
|--------|------|------|
| GET | `/players` | ❌ |

### Palpites
| Método | Rota | Auth |
|--------|------|------|
| GET | `/predictions` | ✅ |
| POST | `/predictions` | ✅ |
| GET | `/predictions/special` | ✅ |
| POST | `/predictions/special` | ✅ |

### Ranking
| Método | Rota | Auth |
|--------|------|------|
| GET | `/ranking` | ❌ |

---

## ⚙️ Como rodar localmente

### Pré-requisitos

- Go 1.22+
- PostgreSQL (ou conta no [Neon](https://neon.tech))
- Conta na [RapidAPI](https://rapidapi.com) com API-Football

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/bolao-copa-2026.git
cd bolao-copa-2026
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
```

Edite o `.env`:

```env
DATABASE_URL=postgres://user:password@host/dbname?sslmode=require
JWT_SECRET=um_segredo_longo_e_aleatorio
API_FOOTBALL_KEY=sua_chave_da_rapidapi
PORT=8080
```

### 3. Instale as dependências

```bash
go mod tidy
```

### 4. Rode o servidor

```bash
go run ./cmd/api
```

As migrations rodam automaticamente na inicialização. O servidor sobe na porta definida em `PORT`.

### 5. Rode os testes

```bash
# Todos os testes
go test ./...

# Com cobertura
go test ./... -cover

# Relatório de cobertura no navegador
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 🏗️ Arquitetura

```
[React/TypeScript Frontend]
        |
        | HTTP/JSON
        v
[Go + Gin — Render]
        |
   +----|----+
   |         |
[PostgreSQL]  [Worker — goroutine]
  (Neon)      time.Ticker 1 min
                   |
             [API-Football]
```

### Princípios aplicados

**Single Responsibility (S)** — cada arquivo tem uma responsabilidade única: `model.go` é domínio, `dto.go` é contrato HTTP, `repository.go` é acesso a dados, `service.go` é lógica de negócio, `handler.go` é HTTP.

**Dependency Inversion (D)** — todas as camadas dependem de interfaces, nunca de implementações concretas. O handler depende de `Service`, o service depende de `Repository`. Isso torna cada camada testável de forma isolada com mocks.

```
Handler → Service (interface)
Service → Repository (interface)
Repository → *sql.DB
```

### Worker

Roda em goroutine separada na mesma instância Go. A cada minuto:

1. Busca jogos ativos na API-Football
2. Atualiza status e placar no banco
3. Ao detectar jogo finalizado, dispara cálculo de pontuação
4. Verifica se fase de grupos encerrou para popular o mata-mata
5. Usa flag `scored` nos palpites para garantir idempotência

---

## 🗄️ Schema

As principais tabelas:

- `users` — participantes do bolão
- `teams` — 48 seleções da Copa
- `players` — candidatos a artilheiro
- `matches` — 104 jogos da Copa
- `predictions` — palpites por jogo (UNIQUE por user+match)
- `special_predictions` — campeão e artilheiro por usuário
- `group_standings` — VIEW calculada dinamicamente (sem tabela separada)

---

## 📄 Documentação

O documento completo de requisitos está em [`docs/SRD.md`](doc.md), cobrindo requisitos funcionais, não funcionais, regras de negócio, schema e critérios de aceite.

---

## 📝 Licença

MIT
