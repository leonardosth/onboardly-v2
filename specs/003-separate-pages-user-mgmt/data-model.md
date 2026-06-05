# Data Model

## Entidades Principais

### 1. Usuário (`users`)
Representa uma conta no sistema. Responsável pelo acesso ao Onboardly.

**Campos Chave**:
- `id` (UUID, Primary Key)
- `email` (String, Unique)
- `password_hash` (String)
- `role` (String, enum: 'Admin', 'Analista')
- `created_at` (Timestamp, padrão ISO 8601)

**Validações Mencionadas na Spec**:
- E-mail único no sistema (FR-006 Edge case)
- Senha deve ter mínimo de 8 caracteres, contendo pelo menos 1 letra e 1 número (FR-006). (Isso será validado no endpoint de criação, antes do hashing).
- Não pode ser excluído se for o último usuário com role 'Admin' (FR-011).

### 2. Projeto (`projects`)
Representa a implantação de um cliente.

**Campos Chave**:
- `id` (UUID, Primary Key)
- `client_id` (UUID, Foreign Key para `clients`)
- `name` (String)
- `status` (String, ex: 'Backlog', 'Em andamento', 'Go-Live')
- `created_at` (Timestamp)

### 3. Cliente (`clients`)
Representa a empresa contratante.

**Campos Chave**:
- `id` (UUID, Primary Key)
- `name` (String)
- `cnpj` (String, Unique)
- `created_at` (Timestamp)
