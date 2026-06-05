# Feature Specification: Separação de Páginas e Gerenciamento de Usuários

**Feature Branch**: `003-separate-pages-user-mgmt`

**Created**: 2026-06-04

**Status**: Draft

**Input**: User description: "Separar Clientes e Projetos em páginas distintas no sidebar e criar tela de gerenciamento de usuários com acesso restrito a administradores."

## Clarifications

### Session 2026-06-04

- Q: O sistema deve impedir a exclusão do último administrador restante? → A: Sim. O sistema DEVE impedir a exclusão de qualquer usuário Admin se ele for o último administrador restante no sistema.
- Q: Quais são os requisitos mínimos de senha para novos usuários? → A: Mínimo de 8 caracteres com pelo menos 1 letra e 1 número.
- Q: A listagem de Projetos deve oferecer filtros ou busca? → A: Sim. Filtro por status (Backlog, Em andamento, Go-Live) e campo de busca por nome do projeto ou cliente.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Navegação separada para Clientes e Projetos (Priority: P1)

Como qualquer usuário autenticado (Admin ou Analista), quero ter itens distintos no menu lateral para "Clientes" e "Projetos", de modo que eu possa acessar cada área de forma independente sem precisar navegar por uma para chegar à outra.

**Why this priority**: A separação é a base da reorganização da navegação. Sem ela, as demais funcionalidades ficam inacessíveis de maneira direta.

**Independent Test**: Pode ser testada ao fazer login e verificar que o sidebar exibe dois links distintos (Clientes e Projetos), cada um levando à respectiva página de listagem.

**Acceptance Scenarios**:

1. **Given** um usuário autenticado no sistema, **When** ele visualiza o menu lateral, **Then** ele vê itens separados para "Clientes" e "Projetos" (além de "Dashboard").
2. **Given** um usuário no Dashboard, **When** ele clica em "Projetos" no sidebar, **Then** ele é levado a uma página de listagem de todos os projetos com informações de nome, cliente associado e status.
3. **Given** um usuário no Dashboard, **When** ele clica em "Clientes" no sidebar, **Then** ele é levado à página de listagem de clientes (que já existe).

---

### User Story 2 - Página dedicada de listagem de Projetos (Priority: P1)

Como qualquer usuário autenticado, quero acessar uma página dedicada que lista todos os projetos de implantação do sistema, exibindo o nome do projeto, o cliente ao qual pertence e o status atual, para ter uma visão centralizada do portfólio de projetos.

**Why this priority**: Sem esta página, o item "Projetos" no sidebar não teria destino. É complementar à User Story 1.

**Independent Test**: Pode ser testada acessando `/projects` e verificando que todos os projetos existentes são listados com os dados corretos.

**Acceptance Scenarios**:

1. **Given** existem projetos cadastrados no sistema, **When** o usuário acessa a página de Projetos, **Then** o sistema exibe uma lista com o nome do projeto, o nome do cliente vinculado e o status atual (ex: "Backlog", "Em andamento", "Go-Live").
2. **Given** não existem projetos cadastrados, **When** o usuário acessa a página de Projetos, **Then** o sistema exibe a mensagem "Nenhum projeto cadastrado."
3. **Given** um usuário está na listagem de Projetos, **When** ele clica em um projeto, **Then** ele é redirecionado para a página de detalhes daquele projeto (que já existe).
4. **Given** existem projetos com status variados, **When** o usuário seleciona um filtro de status (ex: "Em andamento"), **Then** a lista exibe apenas os projetos com o status selecionado.
5. **Given** existem projetos cadastrados, **When** o usuário digita um termo no campo de busca, **Then** a lista é filtrada em tempo real por nome do projeto ou nome do cliente associado.

---

### User Story 3 - Gerenciamento de Usuários restrito a Administradores (Priority: P2)

Como administrador, quero acessar uma tela dedicada para gerenciar os usuários do sistema (listar, criar e excluir), para que eu possa controlar quem tem acesso à plataforma e com qual nível de permissão.

**Why this priority**: Complementa a proposta de segurança RBAC da constituição do projeto, mas depende das funcionalidades de navegação (P1) já estarem prontas.

**Independent Test**: Pode ser testada fazendo login como Admin, acessando a página de Usuários, e verificando que é possível visualizar a lista, cadastrar um novo usuário e excluir um existente. Também verificar que um usuário Analista não consegue acessar essa página.

**Acceptance Scenarios**:

1. **Given** um administrador autenticado, **When** ele visualiza o menu lateral, **Then** ele vê o item "Usuários" no sidebar.
2. **Given** um analista autenticado, **When** ele visualiza o menu lateral, **Then** o item "Usuários" NÃO aparece no sidebar.
3. **Given** um administrador na tela de Usuários, **When** ele visualiza a lista, **Then** o sistema exibe todos os usuários registrados com e-mail, papel (Admin/Analista) e data de criação.
4. **Given** um administrador na tela de Usuários, **When** ele preenche o formulário de novo usuário com e-mail, senha e papel e clica em Salvar, **Then** o sistema cria o novo usuário e ele aparece na lista.
5. **Given** um administrador na tela de Usuários, **When** ele clica em "Excluir" em um usuário, **Then** o sistema solicita confirmação e, se confirmada, remove o usuário da lista.
6. **Given** um analista tenta acessar diretamente a URL `/users`, **When** a página carrega, **Then** o sistema redireciona o analista de volta ao Dashboard com uma mensagem de acesso negado ou simplesmente não renderiza a tela.

---

### Edge Cases

- O que acontece se um administrador tentar excluir a si mesmo? O sistema deve impedir essa ação com uma mensagem como "Você não pode excluir a própria conta."
- O que acontece se o administrador tenta criar um usuário com um e-mail já existente? O sistema deve exibir um erro indicando duplicidade.
- O que acontece se um projeto na listagem de projetos não possuir um cliente associado (dados inconsistentes)? O sistema deve exibir "Cliente não encontrado" naquele campo.
- O que acontece se o administrador tenta excluir o último administrador restante no sistema? O sistema deve impedir a exclusão e exibir uma mensagem como "Não é possível excluir o último administrador do sistema."
- O que acontece se a senha informada ao criar um novo usuário não atender aos requisitos mínimos (8 caracteres, 1 letra, 1 número)? O sistema deve exibir uma mensagem de validação indicando os critérios não atendidos antes de permitir o envio.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: O sistema DEVE exibir no sidebar três itens de navegação independentes: "Dashboard", "Clientes" e "Projetos" para todos os usuários autenticados.
- **FR-002**: O sistema DEVE exibir um quarto item "Usuários" no sidebar, visível exclusivamente para usuários com papel "Admin".
- **FR-003**: O sistema DEVE disponibilizar uma página dedicada `/projects` que lista todos os projetos de implantação, mostrando nome, cliente associado e status. A página DEVE incluir filtro por status (Backlog, Em andamento, Go-Live) e campo de busca por nome do projeto ou nome do cliente.
- **FR-004**: O sistema DEVE disponibilizar uma página dedicada `/users` para o gerenciamento de usuários, acessível apenas por administradores.
- **FR-005**: A tela de Usuários DEVE permitir listar todos os usuários com e-mail, papel e data de criação.
- **FR-006**: A tela de Usuários DEVE permitir a criação de novos usuários com e-mail, senha e papel (Admin ou Analista). A senha DEVE ter no mínimo 8 caracteres, contendo pelo menos 1 letra e 1 número.
- **FR-007**: A tela de Usuários DEVE permitir a exclusão de usuários, com confirmação antes da ação.
- **FR-008**: O sistema DEVE impedir que um administrador exclua a própria conta.
- **FR-009**: O sistema DEVE redirecionar ou bloquear qualquer tentativa de acesso à rota `/users` por usuários não-administradores.
- **FR-010**: A página de listagem de Projetos DEVE exibir uma mensagem de estado vazio ("Nenhum projeto cadastrado.") quando não houver projetos no sistema.
- **FR-011**: O sistema DEVE impedir a exclusão de qualquer usuário com papel Admin quando ele for o último administrador restante, exibindo mensagem de erro apropriada.

### Key Entities

- **Usuário**: Representa uma conta do sistema com e-mail, senha criptografada, papel (Admin ou Analista) e data de criação. Relacionamento: um usuário pode ser Analista de vários projetos/reuniões.
- **Projeto**: Representa um projeto de implantação vinculado a um Cliente. Atributos-chave: nome, status, cliente associado, flag ativo/inativo.
- **Cliente**: Representa a empresa-cliente cadastrada. Atributos-chave: nome, CNPJ, data de criação.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Qualquer usuário autenticado consegue navegar do Dashboard para a lista de Clientes ou de Projetos em no máximo 1 clique a partir do sidebar.
- **SC-002**: Administradores conseguem criar um novo usuário no sistema em menos de 30 segundos (preenchimento e submissão do formulário).
- **SC-003**: Analistas não conseguem acessar ou visualizar o conteúdo da tela de Usuários, mesmo digitando a URL diretamente.
- **SC-004**: A listagem de projetos carrega e exibe 100% dos projetos existentes no banco de dados com seus respectivos clientes e status.

## Assumptions

- O sistema de autenticação JWT e RBAC (Admin/Analista) existente será reutilizado sem modificações.
- A API do backend para listagem de projetos (`GET /api/projects`) já existe e será reaproveitada; uma nova rota será necessária para listagem/exclusão de usuários.
- As páginas de detalhe de projeto (`/projects/:id`) e detalhe de cliente (`/clients/:id`) já existem e continuarão funcionando normalmente.
- O visual e os estilos seguirão o padrão Dark Theme já estabelecido no restante da aplicação (sidebar, modais, cards).
- A exclusão de um usuário não afeta dados históricos de reuniões ou projetos já vinculados a ele (soft reference).
