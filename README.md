# Rinha de Backend 2024 Q1 - Controle de Concorrência

Esta é uma aplicação em Golang desenvolvida com o framework Gin e utiliza o banco de dados PostgreSQL. A aplicação é destinada a participar da [Rinha de Backend 2024 Q1](https://github.com/zanfranceschi/rinha-de-backend-2024-q1)! para controle de concorrência.

## Funcionalidades

A aplicação possui dois endpoints:

1. **POST /clientes/[id]/transacoes**: Este endpoint permite registrar transações para um cliente específico. Requer um corpo JSON com os seguintes campos:

   ```json
   {
     "valor": 1000,
     "tipo": "c",
     "descricao": "descricao"
   }
   ```

   - `valor`: O valor da transação.
   - `tipo`: O tipo de transação (por exemplo, "c" para crédito).
   - `descricao`: Descrição da transação.

2. **GET /clientes/[id]/extrato**: Este endpoint permite obter o extrato de transações de um cliente específico.

## Configuração

1. **Instalação de Dependências**:
   Certifique-se de ter o Go instalado em sua máquina.

   ```bash
   go get -u github.com/gin-gonic/gin
   go get -u github.com/jinzhu/gorm
   go get -u github.com/jinzhu/gorm/dialects/postgres
   ```

2. **Configuração do Banco de Dados**:

   - Crie um banco de dados PostgreSQL.
   - Edite o arquivo `config.go` e atualize as informações do banco de dados conforme necessário.

3. **Executando a Aplicação**:
   Execute o seguinte comando na raiz do projeto:
   ```bash
   go run main.go
   ```

## Exemplos de Uso

### Registrar Transação

```bash
curl -X POST \
  http://localhost:8080/clientes/[id]/transacoes \
  -H 'Content-Type: application/json' \
  -d '{
    "valor": 1000,
    "tipo" : "c",
    "descricao" : "descricao"
}'

## Stack

- [Go 1.21](https://go.dev/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [PostgreSQL](https://www.postgresql.org/)

#### [Linkedin](https://br.linkedin.com/in/isadora-souza?original_referer=https%3A%2F%2Fwww.google.com%2F)

#### [Twitter](https://twitter.com/isadoraamsouza)
```
