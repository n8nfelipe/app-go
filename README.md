# app-go

API REST de gestão de clientes com autenticação JWT, escrita em Go (Gin) e PostgreSQL.

## Índice

- [Sobre](#sobre)
- [Stack](#stack)
- [Arquitetura](#arquitetura)
- [Pré-requisitos](#pré-requisitos)
- [Como rodar](#como-rodar)
- [Variáveis de ambiente](#variáveis-de-ambiente)
- [Endpoints](#endpoints)
- [Estrutura de pastas](#estrutura-de-pastas)
- [Migrations](#migrations)
- [Testes manuais](#testes-manuais)
- [Comandos úteis](#comandos-úteis)

## Sobre

API para cadastro e consulta de clientes, com:

- Autenticação por **JWT (HS256)**, expiração de 2 horas
- Senhas armazenadas com **Argon2id** (com salt aleatório)
- Configuração por **variáveis de ambiente** (Viper + godotenv)
- Migrations SQL para PostgreSQL
- Camadas separadas: `controller → usecase → repository → model`

## Stack

| Camada     | Tecnologia                          |
| ---------- | ----------------------------------- |
| HTTP       | [Gin](https://github.com/gin-gonic/gin) |
| Auth       | `golang-jwt/jwt/v5`                 |
| Senhas     | `golang.org/x/crypto` (Argon2id)    |
| Banco      | PostgreSQL 16 (`lib/pq`, `database/sql`) |
| Config     | Viper + godotenv                    |
| Infra      | Docker Compose                      |

## Arquitetura

```
handle/        → composição das rotas e injeção de dependências
  └─ middleware/  → validação do JWT nas rotas protegidas
internal/
  controller/  → bindings HTTP, status codes, respostas JSON
  usecase/     → regras de negócio (ex.: validar credenciais, hashear senha)
  repositorie/ → queries SQL
  model/       → structs e persistência
  security/    → Argon2id (hash/verify)
config/        → carga de configuração (env + defaults)
db/            → conexão + migrations
cmd/           → pontos de entrada (servidor e CLI de migrations)
```

Fluxo de uma requisição:

```
HTTP → middleware (JWT) → controller → usecase → repository → PostgreSQL
```

## Pré-requisitos

- Go **1.25+**
- Docker + Docker Compose
- (opcional) REST Client do VS Code para usar o `test.http`

## Como rodar

```bash
# 1. Clonar
git clone <repo-url> && cd app-go

# 2. Criar o arquivo de configuração (ajuste senhas/segredos)
cp .env.example .env   # ou crie manualmente — veja a seção de env

# 3. Subir o PostgreSQL
docker compose up -d db

# 4. Aplicar as migrations
go run ./cmd/migrate up

# 5. Subir a API (porta 8081)
go run ./cmd/main.go
```

A API sobe em `http://localhost:8081`.

> O `.env` está no `.gitignore` — **nunca** commite esse arquivo.

## Variáveis de ambiente

Todas usam o prefixo `APP_` (ex.: `db.host` → `APP_DB_HOST`), com defaults entre parênteses.

| Variável            | Descrição                          | Default     |
| ------------------- | ---------------------------------- | ----------- |
| `APP_DB_HOST`       | Host do PostgreSQL                 | `localhost` |
| `APP_DB_PORT`       | Porta do PostgreSQL                | `5432`      |
| `APP_DB_USER`       | Usuário do banco                   | `postgres`  |
| `APP_DB_PASSWORD`   | Senha do banco                     | `password`  |
| `APP_DB_DBNAME`     | Nome do banco                      | `mydb`      |
| `JWT_SECRET`        | Segredo de assinatura do token JWT | *(obrigatório)* |

Exemplo de `.env`:

```env
APP_DB_HOST=localhost
APP_DB_PORT=5432
APP_DB_USER=postgres
APP_DB_PASSWORD=password
APP_DB_DBNAME=mydb
JWT_SECRET=troque-por-um-segredo-forte
```

O `docker-compose.yml` usa as mesmas variáveis (`APP_DB_*`) para criar o container do Postgres.

## Endpoints

### Públicos

| Método | Rota       | Descrição                         |
| ------ | ---------- | --------------------------------- |
| `POST` | `/login`   | Autentica e retorna o token JWT   |
| `POST` | `/users`   | Cria usuário (senha hasheada)     |

### Protegidos (header `Authorization: Bearer <token>`)

| Método | Rota           | Descrição                 |
| ------ | -------------- | ------------------------- |
| `GET`  | `/clients/`    | Lista todos os clientes   |
| `GET`  | `/clients/:id` | Busca cliente por ID      |
| `POST` | `/clients/`    | Cria um cliente           |
| `PUT`  | `/clients/:id` | Atualiza um cliente       |

> Sem token (ou com token expirado) as rotas de clientes respondem `401`.

### Corpos das requisições

**`POST /login`**

```json
{ "username": "felipe", "password": "secret" }
```

Resposta `200`:

```json
{ "token": "eyJhbGciOi...", "userName": "felipe" }
```

**`POST /users`**

```json
{ "username": "felipe", "password": "secret" }
```

**`POST /clients/`** e **`PUT /clients/:id`**

```json
{ "name": "Empresa X", "email": "contato@empresax.com" }
```

Resposta de criação retorna o `id` gerado.

**`GET /clients/:id`** → objeto do cliente:

```json
{ "id": 1, "name": "Empresa X", "email": "contato@empresax.com", "created_at": "2026-10-21T12:00:00Z" }
```

### Exemplo com curl

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8081/login \
  -H "Content-Type: application/json" \
  -d '{"username":"felipe","password":"secret"}' | jq -r .token)

# Criar cliente
curl -X POST http://localhost:8081/clients/ \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Empresa X","email":"contato@empresax.com"}'

# Listar clientes
curl http://localhost:8081/clients/ -H "Authorization: Bearer $TOKEN"
```

## Estrutura de pastas

```
.
├── cmd/
│   ├── main.go            # servidor HTTP (porta 8081)
│   └── migrate/main.go    # CLI: up | down | version
├── config/con.go          # Viper + godotenv
├── db/
│   ├── connection.go      # abertura/ping da conexão
│   ├── migrations/        # .up.sql / .down.sql
│   └── ...                # funções Up/Down/Version
├── handle/
│   ├── client/            # rotas de clientes (grupo protegido)
│   └── user/              # rotas de login e criação de usuário
├── internal/
│   ├── controller/        # HTTP handlers
│   ├── model/             # structs + SQL de persistência
│   ├── repositorie/       # acesso a dados
│   ├── security/          # Argon2id
│   └── usecase/           # regras de negócio
├── middleware/autenticacao.go  # geração e validação do JWT
├── docker-compose.yml     # PostgreSQL 16
├── test.http              # exemplos para REST Client
└── go.mod / go.sum
```

## Migrations

CLI em `cmd/migrate`, com arquivos em `db/migrations/` no formato
`NNN_descricao.up.sql` / `NNN_descricao.down.sql`.

| Comando                   | Descrição                          |
| ------------------------- | ---------------------------------- |
| `go run ./cmd/migrate up`     | Aplica todas as migrations pendentes |
| `go run ./cmd/migrate down`   | Reverte a última migration         |
| `go run ./cmd/migrate version`| Mostra a versão atual (e se está *dirty*) |

Tabelas criadas:

- **clients**: `id`, `name`, `email` (unique), `created_at`
- **users**: `id`, `username` (unique), `password` (hash Argon2id), `created_at`

## Testes manuais

O arquivo [`test.http`](./test.http) contém exemplos para a extensão
**REST Client** do VS Code. Troque o token pelo retornado no `/login`
(tokens expiram em 2 horas).

## Comandos úteis

```bash
go run ./cmd/main.go        # sobe a API
go vet ./...                # análise estática
gofmt -w .                  # formata o código
docker compose up -d db     # sobe só o banco
docker compose down          # para o banco (mantém o volume)
```

## Licença

Distribuído sob a licença [MIT](./LICENSE).
