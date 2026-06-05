# Phase 0: Research & Architecture Decisions

## Decisions

### 1. Backend Route Structure para Usuários
- **Decision**: Criar um novo handler `internal/api/user.go` (ou expandir `auth/handler.go` se fizer sentido, mas preferencialmente isolar o CRUD de usuários em `internal/user/handler.go` e `internal/user/service.go`) protegido pelo middleware de autenticação e por uma verificação de *Role* (apenas "Admin").
- **Rationale**: O princípio I da Constituição exige RBAC estrito. Isolar a lógica garante que endpoints de sistema fiquem longe da lógica de login.
- **Alternatives considered**: Colocar tudo junto com `auth`. Rejeitado por misturar conceitos (Auth = gerenciar sessão, User = CRUD de contas).

### 2. Deleção do Último Administrador
- **Decision**: No `UserService.DeleteUser`, fazer um `COUNT` na base de dados para o role `Admin`. Se o ID a ser deletado for o de um Admin e a contagem de admins for 1, retornar um erro `400 Bad Request`.
- **Rationale**: Previne o estado irrecuperável definido na clarificação Q1 (FR-011).

### 3. Frontend: Projetos e Busca/Filtro
- **Decision**: Criar `ProjectsList.vue` contendo um grid ou tabela de projetos. Utilizar *computed properties* do Vue para filtrar localmente a lista de projetos provinda do estado (Pinia `projectsStore`) por `status` e por `searchQuery` (nome do projeto ou cliente).
- **Rationale**: Implementação simples e reativa, adequada para o volume de dados esperado inicialmente.

### 4. Frontend: Sidebar
- **Decision**: Atualizar `MainLayout.vue` para incluir os links: `/dashboard`, `/clients`, `/projects`, e, condicionalmente se `authStore.isAdmin == true`, `/users`.
- **Rationale**: Mantém o layout unificado e aplica a restrição visual (FR-002).
