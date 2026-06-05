<!--
SYNC IMPACT REPORT:
- Version change: Template -> 1.0.0
- List of modified principles:
  - PRINCIPLE_1: [PRINCIPLE_1_NAME] -> I. Security & Access Control (RBAC)
  - PRINCIPLE_2: [PRINCIPLE_2_NAME] -> II. Referential Integrity & Relational Rules
  - PRINCIPLE_3: [PRINCIPLE_3_NAME] -> III. Standardized Chronological & Date Protocols
  - PRINCIPLE_4: [PRINCIPLE_4_NAME] -> IV. Robust External Integrations
  - PRINCIPLE_5: [PRINCIPLE_5_NAME] -> V. Real-time Consolidated Metrics & KPIs
- Added sections:
  - Logical Modules & Functional Scope (Section 2)
  - Development Workflow & Quality Gates (Section 3)
- Removed sections: None
- Templates requiring updates:
  - ✅ C:/Users/leona/Documents/Projects/onboardly-v2/.specify/templates/plan-template.md
  - ✅ C:/Users/leona/Documents/Projects/onboardly-v2/.specify/templates/spec-template.md
  - ✅ C:/Users/leona/Documents/Projects/onboardly-v2/.specify/templates/tasks-template.md
- Follow-up TODOs: None
-->

# Onboardly Constitution

## Core Principles

### I. Security & Access Control (RBAC)
The system MUST enforce strict Role-Based Access Control (RBAC) dynamically. User credentials MUST
be hashed using bcrypt prior to database storage. Sessions MUST be validated via cryptographically
signed JSON Web Tokens (JWT) for tracing and securing API transactions.

### II. Referential Integrity & Relational Rules
The database MUST enforce referential integrity and strict field-level validations. Agendamentos
de reuniões e interações MUST NOT be created if they reference a non-existent or inactive project.
Cadastro de novos clientes MUST validate essential attributes such as CNPJ and Name before saving.

### III. Standardized Chronological & Date Protocols
All datetime properties within input payloads, API responses, and database schemas MUST follow
the ISO 8601/RFC3339 formatting standards. This guarantees chronological coherence across features
such as dashboard histories, analytics, and meeting schedules.

### IV. Robust External Integrations
External integrations, including ERP synchronization, MUST be handled through robust Webhooks
processing JSON payloads. Handlers MUST validate raw input schemas and log sync statuses for
auditing.

### V. Real-time Consolidated Metrics & KPIs
Dashboard statistics, success/failure KPIs (such as activation rate and meeting no-shows), and
monthly activity feeds MUST be derived dynamically using robust database queries, avoiding arbitrary
estimates, and maintaining a 6-month timeline window.

## Logical Modules & Functional Scope

### 1. Autenticação e Autorização
- **RF01**: Registro de novos usuários com níveis de acesso (ex: Admin, Analista).
- **RF02**: Autenticação por e-mail e senha com hash bcrypt.
- **RF03**: Geração e validação de tokens de sessão (JWT) para rastreabilidade.
- **RF04**: Controle de acesso baseado em papéis (RBAC).

### 2. Gestão de Clientes
- **RF05**: Cadastro de novos clientes com validação de Nome e CNPJ.
- **RF06**: Recebimento de dados cadastrais e contratos via Webhooks/JSON do ERP Corporativo.
- **RF07**: CRUD completo para clientes (listar, visualizar detalhes, atualizar, excluir).

### 3. Projetos de Implantação
- **RF08**: Criação de Projetos de Implantação vinculados a um Cliente.
- **RF09**: Acompanhamento e atualização de status (Backlog, Em andamento, Go-Live, etc.).
- **RF10**: Listagem de projetos ativos da carteira com filtros por analista ou status.
- **RF11**: Cálculo ou manutenção do status de ativação (ativo/inativo) do projeto.

### 4. Reuniões e Interações
- **RF12**: Agendamento de reuniões atreladas a projetos específicos.
- **RF13**: Registro do Analista responsável que conduz a reunião.
- **RF14**: Garantia de integridade referencial (impedir reunião em projeto inexistente).
- **RF15**: Padronização de datas no formato ISO 8601/RFC3339.

### 5. Dashboard e Analytics (BI)
- **RF16**: Painel principal (Dashboard) com métricas consolidadas em tempo real.
- **RF17**: Cálculo de taxas de sucesso/falha (índice de ativação ideal ~86%, no-show <10%).
- **RF18**: Histórico mensal das implantações nos últimos 6 meses.
- **RF19**: Feed de "Atividades Recentes" unindo clientes, projetos e reuniões.

### 6. Notificações e Relatórios
- **RF20**: Emissão de notificações ou alertas visuais para prazos de projetos próximos do vencimento.
- **RF21**: Exportação de relatórios de produtividade em formato PDF.

## Development Workflow & Quality Gates

- **Testing Gate**: Any core business logic (auth, client CRUD, status calculation) MUST be verified by integration and unit tests when explicitly required in specifications.
- **Migration Schema**: Database schema modifications MUST be done via structured migration files. Directly editing the active database structure is forbidden.
- **Branch Management**: Feature implementations MUST follow branch naming policies using the format "###-feature-name" (where ### is the sequential or timestamp ID).

## Governance

1. **Constitutional Authority**: This constitution defines the core architectural guidelines and requirements of Onboardly. Any feature implementation MUST adhere to these rules.
2. **Complexity Tracking**: Any design decision that violates a core principle MUST be documented and justified in the plan's complexity tracking table.
3. **Amendments**: Modifying this constitution requires updating the version number in accordance with semantic versioning rules and prepending a Sync Impact Report.

**Version**: 1.0.0 | **Ratified**: 2026-06-04 | **Last Amended**: 2026-06-04
