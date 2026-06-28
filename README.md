# Onboardly V2

[![CI Pipeline](https://github.com/leonardosth/onboardly-v2/actions/workflows/ci.yml/badge.svg)](https://github.com/leonardosth/onboardly-v2/actions/workflows/ci.yml)

**Onboardly V2** é uma plataforma moderna e interativa focada no gerenciamento de implantações (deployments) e no processo de integração de clientes (onboarding). O sistema permite acompanhar o status de projetos, gerenciar a base de clientes, agendar reuniões (agendas) com interface de calendário rica e monitorar os KPIs e métricas da operação em tempo real através de um painel visual avançado.

---

## 🎯 Funcionalidades Principais

- **Dashboard de Métricas**: Visualização em tempo real das métricas da operação, com indicadores-chave (KPIs) sobre clientes ativos, agendas pendentes e status geral da integração. Conta também com um painel de atividades recentes (Audit Log).
- **Gestão de Clientes e Projetos**: Fluxo completo (CRUD) para cadastrar clientes e vincular projetos a eles. Os projetos avançam por diferentes fases operacionais: *Backlog*, *Em andamento* e *Go-Live*.
- **Agendamentos e Calendário interativo**: Sistema para organizar as reuniões de onboarding. Possui integração completa com interface de Calendário (via FullCalendar), oferecendo visões mensais e semanais, indicativos de cores e facilidade de manipulação.
- **Integração com Webhooks**: Possui endpoint dedicado (protegido por token) para sincronizar e importar clientes externamente — por exemplo, oriundos de integrações via Google Sheets, Make ou Zapier.
- **Autenticação e Permissões**: Login via JWT. Perfil Administrativo (visão geral, gerenciamento de usuários, exclusão de dados) e perfil de Analista (foco nas agendas e evolução técnica dos projetos).

---

## 🏗️ Arquitetura

O sistema adota uma **Arquitetura Cliente-Servidor (Monolito modularizado no backend)**, visando performance, tipagem forte e rapidez no desenvolvimento:

- **Backend (API RESTful)**: Construído em Go (Golang) para garantir latência ultrabaixa e controle fino sobre concorrência e memória. A aplicação segue os padrões da comunidade Go: 
  - *Handlers*: Controladores que interagem com o HTTP (`go-chi`).
  - *Services / Repositories*: Onde as lógicas de negócios operam diretamente interagindo com o banco via comandos raw SQL otimizados (`pgx`).
- **Frontend (SPA - Single Page Application)**: Arquitetado em Vue.js 3 usando a Composition API. A reatividade do estado global da interface é delegada ao `Pinia`, a transição entre telas ao `Vue Router`, e as requisições à API tratadas por intermédio do `Axios`. Todo o ecossistema é servido pelo `Vite`.
- **Database (PostgreSQL)**: Banco de dados relacional sólido para garantir a integridade referencial dos dados entre Usuários, Clientes, Projetos e Reuniões. Conta com restrições (`Foreign Keys`) e tabelas projetadas para análises analíticas (Dashboard).

---

## 🛠️ Tecnologias e Stacks

### Backend
- **Linguagem**: Go 1.22+
- **Roteador HTTP**: [go-chi/chi](https://github.com/go-chi/chi) (Leve, nativamente idiomático, middlewares embutidos).
- **Database Driver**: [jackc/pgx](https://github.com/jackc/pgx) (Alta performance para Postgres em Go).
- **Segurança**: bcrypt para senhas e `golang-jwt` para gerenciamento das sessões stateless.
- **Cors**: `go-chi/cors`.
- **Testes**: `testing` (stdlib) + [testify](https://github.com/stretchr/testify) (assertions) + [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) (mock de banco).
- **Linting**: [golangci-lint](https://golangci-lint.run/) (`errcheck`, `govet`, `staticcheck`, `gofmt`, `goimports`).

### Frontend
- **Core**: Vue.js 3 (Composition API) + Vite.
- **State Management**: Pinia.
- **Componentes Extra**: FullCalendar (`@fullcalendar/vue3`) para renderização fluida das agendas.
- **Estilização**: Vanilla CSS com temática Cyberpunk/Glassmorphism (Fundos noturnos, `backdrop-filter`, paleta vibrante de cianos e esmeraldas).
- **Linting**: ESLint + `eslint-plugin-vue`.

### DevOps
- **CI/CD**: GitHub Actions (testes, lint, build Docker).
- **Containerização**: Docker multi-stage builds + Docker Compose.
- **Qualidade**: SonarQube (análise estática local).

---

## 🚀 Setup & Instalação

### Requisitos Prévios
- Go 1.22+
- Node.js 18+ (recomendado Node 20+)
- PostgreSQL rodando local ou remotamente.

### Opção 1: Setup Manual (Desenvolvimento)

#### Passo 1: Configuração do Backend
1. Navegue para o diretório `backend/` e instale as dependências:
   ```bash
   cd backend
   go mod tidy
   ```
2. Crie e configure o arquivo de variáveis de ambiente `.env`:
   ```bash
   cp .env.example .env
   ```
   Ajuste as credenciais, principalmente `DATABASE_URL` e `JWT_SECRET`.
3. Inicie o banco de dados rodando as migrations em `backend/migrations/`:
   ```bash
   psql -U postgres -d onboardlyv2 -f migrations/000001_init.up.sql
   psql -U postgres -d onboardlyv2 -f migrations/000002_activation.up.sql
   ```
4. Inicie o servidor (O backend estará ouvindo por padrão em `http://localhost:8080`):
   ```bash
   go run cmd/server/main.go
   ```

#### Passo 2: Configuração do Frontend
1. Navegue para a pasta `frontend/`:
   ```bash
   cd frontend
   npm install
   ```
2. Suba o servidor de desenvolvimento do Vite:
   ```bash
   npm run dev
   ```
3. Acesse a aplicação no seu navegador: `http://localhost:5173`.

### Opção 2: Docker Compose (Produção / Ambiente Completo)

Suba todos os serviços com um único comando:

```bash
# Configure as variáveis de ambiente (opcional, tem defaults)
export DB_PASSWORD=sua_senha_segura
export JWT_SECRET=seu_jwt_secret
export WEBHOOK_TOKEN=seu_webhook_token

# Inicie todos os serviços
docker compose up --build -d
```

Acesse a aplicação em `http://localhost`.

| Serviço    | Porta | Descrição                        |
|------------|-------|----------------------------------|
| `frontend` | 80    | SPA Vue.js via Nginx             |
| `backend`  | 8080  | API RESTful Go                   |
| `db`       | 5432  | PostgreSQL 15                    |

Para parar os serviços:
```bash
docker compose down
```

---

## 🧪 Testes

### Backend (Go)
```bash
cd backend

# Rodar todos os testes com output detalhado
go test ./internal/... -v

# Rodar com cobertura
go test ./internal/... -v -coverprofile=coverage.out

# Verificar análise estática
go vet ./internal/... ./cmd/...
```

#### Cobertura de Testes

| Pacote       | Arquivo(s) de Teste                    | Cenários |
|--------------|----------------------------------------|----------|
| `apierr`     | `errors_test.go`                       | Status codes, JSON body, Content-Type |
| `activity`   | `service_test.go`                      | LogActivity, GetRecentActivities, DB errors |
| `auth`       | `service_test.go`                      | Hash, CheckPassword, JWT (claims, expiração, secret inválido), CreateUser |
| `client`     | `model_test.go`, `service_test.go`     | Validação CNPJ, GetClientsWithDetails, DB errors |
| `config`     | `config_test.go`                       | Defaults, overrides via env vars |
| `dashboard`  | `service_test.go`                      | Métricas completas, divisão por zero |
| `meeting`    | `service_test.go`                      | CreateMeeting, CompleteMeeting, GetByProject/Analyst |
| `project`    | `service_test.go`                      | Status validation, nome vazio |
| `user`       | `service_test.go`, `handler_test.go`   | DeleteUser (self, last admin), ListUsers, validatePassword |

### Frontend (ESLint)
```bash
cd frontend

# Rodar linting
npm run lint
```

---

## 🔄 CI/CD Pipeline

O projeto utiliza **GitHub Actions** para integração contínua. O pipeline é acionado automaticamente em **push** e **pull requests** na branch `main`.

### Jobs do Pipeline

| Job                    | Descrição                                          |
|------------------------|----------------------------------------------------|
| `backend-tests`        | Go vet + testes unitários com cobertura             |
| `frontend-lint-build`  | ESLint + build de produção do Vite                 |
| `docker-build`         | Validação do build das imagens Docker              |

O arquivo de configuração está em `.github/workflows/ci.yml`.

---

## 📡 API Endpoints

### Autenticação (Pública)

| Método | Rota                  | Descrição                     |
|--------|-----------------------|-------------------------------|
| POST   | `/api/auth/register`  | Criar conta de usuário        |
| POST   | `/api/auth/login`     | Autenticar e obter JWT token  |

### Webhooks (Token)

| Método | Rota                     | Descrição                              |
|--------|--------------------------|----------------------------------------|
| POST   | `/api/webhooks/clients`  | Sincronizar clientes via webhook       |

### Clientes (Autenticado)

| Método | Rota                 | Descrição                      | Permissão       |
|--------|----------------------|--------------------------------|-----------------|
| GET    | `/api/clients`       | Listar clientes com detalhes   | Todos           |
| GET    | `/api/clients/{id}`  | Obter cliente por ID           | Todos           |
| POST   | `/api/clients`       | Criar novo cliente             | Todos           |
| PUT    | `/api/clients/{id}`  | Atualizar cliente              | Todos           |
| DELETE | `/api/clients/{id}`  | Excluir cliente                | Admin           |

### Projetos (Autenticado)

| Método | Rota                          | Descrição                     |
|--------|-------------------------------|-------------------------------|
| GET    | `/api/projects`               | Listar projetos               |
| GET    | `/api/projects/{id}`          | Obter projeto por ID          |
| POST   | `/api/projects`               | Criar projeto                 |
| PUT    | `/api/projects/{id}`          | Atualizar status do projeto   |
| POST   | `/api/projects/{id}/finalize` | Finalizar projeto (Go-Live)   |

### Agendas / Reuniões (Autenticado)

| Método | Rota                          | Descrição                       |
|--------|-------------------------------|---------------------------------|
| GET    | `/api/meetings`               | Listar reuniões por projeto     |
| GET    | `/api/meetings/mine`          | Listar minhas reuniões          |
| POST   | `/api/meetings`               | Criar reunião                   |
| POST   | `/api/meetings/{id}/complete` | Marcar reunião como concluída   |

### Usuários (Admin)

| Método | Rota               | Descrição            |
|--------|---------------------|----------------------|
| GET    | `/api/users`        | Listar usuários      |
| POST   | `/api/users`        | Criar usuário        |
| DELETE | `/api/users/{id}`   | Excluir usuário      |

### Dashboard (Autenticado)

| Método | Rota              | Descrição                         |
|--------|--------------------|-----------------------------------|
| GET    | `/api/dashboard`   | Obter métricas e KPIs agregados   |

### Health Check (Público)

| Método | Rota       | Descrição          |
|--------|------------|--------------------|
| GET    | `/health`  | Status do servidor |

---

## 🔐 Conceitos e Modelos de Domínio

- **Usuários (Users)**: Autenticam-se e manipulam o sistema. Roles: `Admin` e `Analista`.
- **Clientes (Clients)**: Pessoas jurídicas ou contas que compraram a solução. Identificados por CNPJ.
- **Projetos (Projects)**: A jornada do cliente dentro do Onboardly. Todo cliente possui ao menos um projeto atrelado que mapeia seu processo de implantação. Status: `Backlog` → `Em andamento` → `Go-Live`.
- **Agendas (Meetings)**: Pontos de sincronia entre o analista e o cliente para evoluir um projeto. Status: `scheduled` → `completed` / `cancelled`.
- **Atividades (Activities)**: Log de auditoria para rastreamento das ações realizadas no sistema.

---

## 📁 Estrutura do Projeto

```
onboardly-v2/
├── .github/workflows/       # Pipeline CI/CD (GitHub Actions)
│   └── ci.yml
├── backend/
│   ├── cmd/
│   │   ├── server/           # Entrypoint da aplicação
│   │   └── seed/             # Script de seed do banco
│   ├── internal/
│   │   ├── activity/         # Serviço de log de atividades
│   │   ├── api/              # Router e middlewares
│   │   ├── apierr/           # Helper de erros HTTP
│   │   ├── auth/             # Autenticação (JWT, bcrypt)
│   │   ├── client/           # CRUD de clientes + webhook
│   │   ├── config/           # Configuração via env vars
│   │   ├── dashboard/        # Métricas e KPIs agregados
│   │   ├── db/               # Pool de conexão PostgreSQL
│   │   ├── meeting/          # Gerenciamento de reuniões
│   │   ├── project/          # Gestão de projetos
│   │   └── user/             # Administração de usuários
│   ├── migrations/           # SQL migrations
│   ├── Dockerfile            # Multi-stage build (Go)
│   ├── .golangci.yml         # Configuração do linter
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/       # Componentes Vue reutilizáveis
│   │   ├── layouts/          # Layouts base
│   │   ├── pages/            # Páginas da aplicação
│   │   ├── router/           # Rotas do Vue Router
│   │   ├── services/         # Serviços Axios
│   │   └── stores/           # Pinia stores
│   ├── Dockerfile            # Multi-stage build (Node + Nginx)
│   ├── nginx.conf            # Config Nginx para SPA
│   └── package.json
├── docker-compose.yml        # Orquestração dos serviços
└── README.md
```
