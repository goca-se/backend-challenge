# EntregasRápidas API Mock

Esta é uma implementação mock da API da transportadora EntregasRápidas para o desafio de gateway de frete.

## Características

- Formato XML para requisições e respostas
- Autenticação via Basic Auth (usuário: entregasrapidas, senha: Senha@123)
- Simula alta latência (2-5 segundos)
- Implementado com FastAPI e Python

## Executando a API

### Com Docker Compose

```bash
docker-compose up -d
```

A API estará disponível em: http://localhost:3000

### Sem Docker (Desenvolvimento)

1. Instale as dependências:

```bash
pip install -r requirements.txt
```

2. Execute a aplicação:

```bash
uvicorn app.main:app --host 0.0.0.0 --port 3000 --reload
```

## Endpoints

### Cotação de Frete

**Endpoint:** `POST /api/v1/entregas-rapidas/cotacao`

**Headers necessários:**
```
Content-Type: application/xml
Authorization: Basic ZW50cmVnYXM6cmFwaWRhcw==
```

**Exemplo de payload:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<cotacao>
  <remetente>
    <cep>01310100</cep>
  </remetente>
  <destinatario>
    <cep>04538132</cep>
  </destinatario>
  <pacote>
    <peso>1.5</peso>
    <altura>20</altura>
    <largura>30</largura>
    <comprimento>15</comprimento>
    <valor_declarado>150.00</valor_declarado>
  </pacote>
  <servico>
    <tipo>todos</tipo>
  </servico>
</cotacao>
```

### Health Check

**Endpoint:** `GET /health`

Retorna o status da API.
