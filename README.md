# Onboardly V2

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

### Frontend
- **Core**: Vue.js 3 (Composition API) + Vite.
- **State Management**: Pinia.
- **Componentes Extra**: FullCalendar (`@fullcalendar/vue3`) para renderização fluida das agendas.
- **Estilização**: Vanilla CSS com temática Cyberpunk/Glassmorphism (Fundos noturnos, `backdrop-filter`, paleta vibrante de cianos e esmeraldas).

---

## 🚀 Setup & Instalação

### Requisitos Prévios
- Go 1.22+
- Node.js 18+ (recomendado Node 20+)
- PostgreSQL rodando local ou remotamente.

### Passo 1: Configuração do Backend
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
3. Inicie o banco de dados rodando a migration localizada em `backend/migrations/000001_init.up.sql`.
4. Inicie o servidor (O backend estará ouvindo por padrão em `http://localhost:8080`):
   ```bash
   go run cmd/server/main.go
   ```

### Passo 2: Configuração do Frontend
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

---

## 🔐 Conceitos e Modelos de Domínio

- **Usuários (Users)**: Autenticam-se e manipulam o sistema.
- **Clientes (Clients)**: Pessoas jurídicas ou contas que compraram a solução.
- **Projetos (Projects)**: A jornada do cliente dentro do Onboardly. Todo cliente possui ao menos um projeto atrelado que mapeia seu processo de implantação.
- **Agendas (Meetings)**: Pontos de sincronia entre o analista e o cliente para evoluir um projeto do `Backlog` até o `Go-Live`.
