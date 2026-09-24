# GoBet Lite 🚀

Backend de iGaming de alta performance desenvolvido em **Go**, projetado com foco em concorrência, integridade financeira e atualizações em tempo real.

## 🛠️ Tecnologias e Arquitetura

* **Linguagem:** Go (Golang)
* **Injeção de Dependências & Ciclo de Vida:** Uber Fx
* **Banco de Dados:** PostgreSQL com `pgx/v5` (Transações atômicas com bloqueio de linha `FOR UPDATE`)
* **Tempo Real:** WebSockets (`gorilla/websocket`) com gerenciador centralizado (`Hub`)
* **Arquitetura:** Clean Architecture / Modular

---

## ⚙️ Principais Funcionalidades

1. **Gestão de Contas e Usuários:** Cadastro e controle de status de usuários.
2. **Processamento de Apostas à Prova de Race Conditions:** 
   - Utiliza transações isoladas do Postgres.
   - O comando `SELECT ... FOR UPDATE` garante que operações simultâneas na mesma conta não gerem inconsistências de saldo.
3. **Notificações Reativas via WebSocket:**
   - Assim que uma aposta é confirmada no banco, o backend empurra instantaneamente o resultado e o novo saldo atualizado direto para a tela do usuário conectado (`/ws?user_id=X`).

---

## 🚀 Como Executar o Projeto

1. **Clone o repositório e configure as variáveis de ambiente** do banco de dados PostgreSQL.
2. **Execute a aplicação:**
   ```bash
   docker compose up -d