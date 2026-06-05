# API Contracts

## Users API

Todos os endpoints abaixo DEVEM exigir token JWT válido.

### `GET /api/users`
Lista todos os usuários do sistema.
**Autorização**: Apenas `Admin`.
**Response (200 OK)**:
```json
[
  {
    "id": "uuid",
    "email": "user@example.com",
    "role": "Admin",
    "created_at": "2026-06-04T12:00:00Z"
  }
]
```

### `POST /api/users`
Cria um novo usuário.
**Autorização**: Apenas `Admin`.
**Request Payload**:
```json
{
  "email": "novo@exemplo.com",
  "password": "SenhaSegura1",
  "role": "Analista"
}
```
**Response (201 Created)**: Retorna o objeto usuário criado (sem a senha).
**Response (400 Bad Request)**: 
```json
{ "error": "A senha deve conter no mínimo 8 caracteres, 1 letra e 1 número" }
```
**Response (409 Conflict)**:
```json
{ "error": "Email já cadastrado" }
```

### `DELETE /api/users/:id`
Exclui um usuário.
**Autorização**: Apenas `Admin`.
**Response (204 No Content)**: Deletado com sucesso.
**Response (400 Bad Request)**:
```json
{ "error": "Não é possível excluir o próprio usuário." }
```
ou
```json
{ "error": "Não é possível excluir o último administrador do sistema." }
```

## Projects API (Existente / Adaptação)

### `GET /api/projects`
Lista todos os projetos do sistema com seus respectivos clientes populados para exibir na tela de listagem de projetos.
**Autorização**: `Admin` ou `Analista`.
**Response (200 OK)**:
```json
[
  {
    "id": "uuid",
    "name": "Implantação ERP",
    "status": "Em andamento",
    "client": {
      "id": "uuid",
      "name": "Empresa XPTO"
    }
  }
]
```
*(Nota: O filtro por status e busca por texto serão aplicados no frontend via Vue computed properties, conforme decidido no Research, então o endpoint só precisa retornar a lista completa).*
