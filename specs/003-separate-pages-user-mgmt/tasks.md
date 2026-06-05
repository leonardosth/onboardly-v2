---
description: "Task list for feature implementation"
---

# Tasks: Separação de Páginas e Gerenciamento de Usuários

**Input**: Design documents from `/specs/003-separate-pages-user-mgmt/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/users-api.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `backend/` e `frontend/src/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

N/A - O projeto já está inicializado.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

N/A - Todas as US desta feature dependem unicamente da estrutura já consolidada no branch principal.

---

## Phase 3: User Story 1 - Navegação separada para Clientes e Projetos (Priority: P1) 🎯 MVP

**Goal**: Ter itens distintos no menu lateral para "Clientes" e "Projetos" acessíveis por Admins e Analistas.

**Independent Test**: Fazer login e verificar que o sidebar exibe os links.

### Implementation for User Story 1

- [x] T001 [P] [US1] Atualizar `frontend/src/layouts/MainLayout.vue` com links do Sidebar para Clientes, Projetos e condicionalmente (v-if="auth.isAdmin") Usuários.
- [x] T002 [P] [US1] Atualizar `frontend/src/router/index.js` adicionando uma rota stub (temporária) para `/projects` apenas para não quebrar a navegação no clique.

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Página dedicada de listagem de Projetos (Priority: P1)

**Goal**: Uma página `/projects` listando projetos com busca textual e filtro de status.

**Independent Test**: Acessar `/projects` via sidebar e testar a exibição da lista e os filtros reativos usando dados reais da API.

### Implementation for User Story 2

- [x] T003 [P] [US2] Criar store do Pinia em `frontend/src/stores/projects.js` para consumir a API `GET /api/projects`.
- [x] T004 [US2] Criar componente `frontend/src/pages/ProjectsList.vue` com o grid/tabela de projetos e computed properties (`searchQuery` e `statusFilter`) baseadas na store de T003.
- [x] T005 [US2] Atualizar `frontend/src/router/index.js` trocando a rota stub do US1 pelo componente real `ProjectsList.vue`.

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Gerenciamento de Usuários restrito a Administradores (Priority: P2)

**Goal**: Tela dedicada para Administradores listarem, criarem (com senha forte) e excluírem usuários (impedindo auto-exclusão e exclusão do último admin).

**Independent Test**: Fazer login como Admin, acessar `/users`, testar CRUD; e logar como Analista para confirmar que o acesso a `/users` é redirecionado/bloqueado.

### Implementation for User Story 3

- [x] T006 [P] [US3] Implementar funções do backend de CRUD de usuário em `backend/internal/user/service.go`, incluindo verificação do último Admin no DB.
- [x] T007 [US3] Criar os handlers HTTP (GET, POST, DELETE) em `backend/internal/user/handler.go` chamando o service de T006 e aplicando validações do payload (regex de senha 8 chars, letras/nums).
- [x] T008 [US3] Registrar rotas `/api/users` no `backend/internal/api/router.go` exigindo o middleware de autorização Admin.
- [x] T009 [P] [US3] Criar store de usuários no frontend em `frontend/src/stores/users.js` mapeando endpoints de T008.
- [x] T010 [US3] Criar `frontend/src/pages/UsersList.vue` incluindo o formulário de cadastro (com tratativa de erro de senha/duplicidade) e o grid de exclusão chamando a store de T009.
- [x] T011 [US3] Atualizar `frontend/src/router/index.js` mapeando a rota `/users` para `UsersList.vue` com `meta: { requiresAuth: true, adminOnly: true }` e protegê-la no `beforeEach`.

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T012 Executar roteiro do `quickstart.md` para validar visualmente todas as rotas de ponta a ponta.
- [x] T013 Remover eventuais lixos e comentários deixados durante o desenvolvimento.

---

## Dependencies & Execution Order

### User Story Dependencies

- **User Story 1 (P1)**: Sem dependências, inicia o MVP.
- **User Story 2 (P1)**: Depende do US1 para se conectar via sidebar.
- **User Story 3 (P2)**: Depende apenas da liberação do backend (podendo rodar paralela ao frontend das US1 e US2).

### Parallel Opportunities

- Todos os endpoints backend (T006, T007, T008) da US3 podem ser iniciados paralelamente ao T001/T002 do frontend.
- O Pinia Store (T003 e T009) pode ser construído antes da componentização (T004 e T010).

---

## Implementation Strategy

### Incremental Delivery

1. O MVP imediato consiste apenas de alterar o Layout com os botões. (US1)
2. Após o US1, constrói-se o consumo da API já existente de Projetos. (US2)
3. Após isso, toda a esteira do US3 (back -> front) é aplicada sem conflitar com nada do sistema atual.
