# Implementation Plan: Separação de Páginas e Gerenciamento de Usuários

**Branch**: `003-separate-pages-user-mgmt` | **Date**: 2026-06-04 | **Spec**: [spec.md](file:///c:/Users/leona/Documents/Projects/onboardly-v2/specs/003-separate-pages-user-mgmt/spec.md)

**Input**: Feature specification from `/specs/003-separate-pages-user-mgmt/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

O objetivo principal desta feature é reorganizar a navegação do sistema e introduzir o controle de usuários. O menu lateral ("Sidebar") será expandido para separar a visão de "Clientes" e "Projetos". Além disso, será introduzida uma tela dedicada e restrita ao Administrador para o gerenciamento de contas de usuário (CRUD), suportada por novos endpoints na API Go.

## Technical Context

**Language/Version**: Go 1.21+ / Vue.js 3

**Primary Dependencies**: Gin Web Framework (Go), GORM (se utilizado no backend) ou pgx, Pinia, Vue Router (Frontend).

**Storage**: PostgreSQL

**Testing**: Dependendo do framework atual (provavelmente pacotes Go nativos e Vitest para o Vue).

**Target Platform**: Web Browser

**Project Type**: Web Application

**Performance Goals**: <200ms API response time para CRUD de usuários e listagem de projetos.

**Constraints**: Acesso baseado em roles rigoroso (JWT Middleware). Validação de senhas no backend. Não deletar o próprio usuário ou o último admin.

**Scale/Scope**: Pequeno a Médio. Novas rotas, 1 novo controller/service no backend, 2 novas telas Vue (ProjectsList, UsersList) e atualização do Layout.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Princípio I (Security & Access Control)**: PASS. O plano define que o novo CRUD de usuários é estritamente validado por roles e que senhas são criadas com segurança.
- **Princípio II (Referential Integrity)**: PASS. A exclusão de usuários não vai apagar projetos para manter integridade, e não é permitido apagar o último Admin.
- **Testing Gate**: PASS. Os fluxos de falha de exclusão podem ser unitariamente testados.

## Project Structure

### Documentation (this feature)

```text
specs/003-separate-pages-user-mgmt/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
backend/
├── internal/
│   ├── api/
│   │   ├── router.go         # Registro das novas rotas de users e projects (se necessário)
│   │   ├── user.go           # Handler para CRUD de usuários
│   │   └── meeting.go        # Handlers existentes
│   ├── auth/                 # Lógica JWT existente
│   └── user/
│       └── service.go        # Regras de negócio de usuário (validação de senha, bloqueio de exclusão do último admin)

frontend/
├── src/
│   ├── layouts/
│   │   └── MainLayout.vue    # Atualizar itens do sidebar (Adicionar Projetos e Usuários, aplicar v-if admin)
│   ├── pages/
│   │   ├── ProjectsList.vue  # Nova tela com cards ou tabela e filtros computados
│   │   └── UsersList.vue     # Nova tela de CRUD de usuários restrita a admin
│   ├── router/
│   │   └── index.js          # Novas rotas com meta: { requiresAdmin: true } para users
│   └── stores/
│       ├── projects.js       # Action para fetch projetos
│       └── users.js          # Nova store para gerenciar usuários
```

**Structure Decision**: A aplicação segue o modelo Option 2 (Web application separada em backend e frontend). O novo código vai focar em `backend/internal/user` e `backend/internal/api` e no frontend criaremos componentes Vue e Pinia stores alinhadas.

## Complexity Tracking

Nenhuma violação constitucional detectada. O design está contido e adere às regras.
