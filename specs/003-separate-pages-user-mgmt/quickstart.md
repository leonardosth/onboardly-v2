# Quickstart Validation Guide

Este documento guia o usuário ou QA na validação ponta a ponta da Feature 003.

## Pré-requisitos
1. Servidor Backend Go rodando (`npm run server` ou `go run cmd/api/main.go`).
2. Frontend Vue.js rodando (`npm run dev`).
3. Banco de dados PostgreSQL rodando e com as migrações aplicadas.
4. Você deve saber a senha da conta de Administrador (ex: `admin@onboardly.com`).

## Cenário 1: Navegação no Sidebar
1. Acesse o sistema e faça login como Administrador ou Analista.
2. Observe o **Menu Lateral**.
3. **Validação**: Deve haver botões visíveis para "Clientes" e "Projetos" separadamente.
4. Clique em "Projetos". A página de listagem de projetos deve carregar mostrando uma lista ou a mensagem "Nenhum projeto cadastrado."

## Cenário 2: Filtros de Projeto
1. Na tela de "Projetos", certifique-se de que há projetos cadastrados.
2. Use a caixa de busca e digite parte do nome de um projeto.
3. **Validação**: A lista deve filtrar instantaneamente os projetos correspondentes.
4. Use o filtro de Select para escolher o status "Go-Live".
5. **Validação**: A lista deve mostrar apenas projetos com esse status.

## Cenário 3: Gerenciamento de Usuários (Apenas Admin)
1. Certifique-se de estar logado como **Admin**.
2. Clique no menu "Usuários" no Sidebar.
3. **Validação**: A lista de usuários do sistema deve aparecer.
4. Clique em "Novo Usuário" e preencha o formulário:
   - Email: `teste@onboardly.com`
   - Senha: `fraca`
   - Role: `Analista`
5. **Validação**: O sistema deve impedir a criação e exibir o aviso "A senha deve conter no mínimo 8 caracteres, 1 letra e 1 número."
6. Corrija a senha para `SenhaForte123` e salve.
7. **Validação**: O novo usuário deve aparecer na lista instantaneamente.

## Cenário 4: Exclusão Segura
1. Tente clicar em "Excluir" no seu próprio usuário (que você usou para fazer o login).
2. **Validação**: O sistema exibe o aviso de que não é possível excluir a própria conta.
3. Clique em "Excluir" no usuário Analista criado no passo anterior e confirme.
4. **Validação**: O usuário desaparece da lista.

## Cenário 5: Bloqueio de Acesso para Analistas
1. Faça logout da conta Admin.
2. Faça login com uma conta de Analista.
3. **Validação**: O menu "Usuários" não deve existir na barra lateral.
4. Tente acessar a URL `http://localhost:5173/users` manualmente no navegador.
5. **Validação**: A página não carrega o conteúdo de usuários ou redireciona para o Dashboard.
