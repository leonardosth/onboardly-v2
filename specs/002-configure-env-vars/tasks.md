# Tasks: Bug Fixes para Renderização do Vue

## Resumo dos Erros
Os componentes Vue (`ClientsList` e `Dashboard`) estão sofrendo crashes ("Unhandled error during execution of render function") devido a tentativas de ler propriedades (`.length` e `.toLowerCase()`) de objetos que estão `null` (retornados pela API ou pelo estado inicial vazio).

## Phase 1: Correções no Frontend

- [x] T001 [BugFix] Atualizar `frontend/src/stores/clients.js` para garantir que `this.clients` sempre receba um array, mesmo se a API retornar `null` (ex: `this.clients = await clientService.getClients() || [];`).
- [x] T002 [BugFix] Atualizar `frontend/src/pages/ClientsList.vue` para utilizar optional chaining ao verificar o tamanho da lista (ex: `clientsStore.clients?.length === 0`).
- [x] T003 [BugFix] Atualizar `frontend/src/pages/Dashboard.vue` para utilizar optional chaining na renderização da badge de role (ex: `:class="auth.role?.toLowerCase() || ''"`).

## Phase 2: Verificação

- [x] T004 [Teste] Acessar o sistema, navegar até `/clients` e garantir que a página carrega corretamente sem erros no console, exibindo a tela vazia caso não haja clientes.
