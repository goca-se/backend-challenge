# Exemplos de APIs das Transportadoras

Este documento contém exemplos de payloads e endpoints para as três transportadoras mockadas que serão utilizadas no desafio.

## Transportadora A - EntregasRápidas

**Características:**
- Formato XML
- Alta latência (2-5 segundos)
- Autenticação via Basic Auth
- Relativamente confiável, mas lenta

### Endpoint
```
POST http://localhost:3000/api/v1/entregas-rapidas/cotacao
```

### Headers
```
Content-Type: application/xml
Authorization: Basic ZW50cmVnYXNyYXBpZGFzOlNlbmhhQDEyMw==
```

### Payload da Requisição
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

### Resposta de Sucesso
```xml
<?xml version="1.0" encoding="UTF-8"?>
<resultado>
  <status>sucesso</status>
  <opcoes>
    <opcao>
      <tipo>normal</tipo>
      <codigo>ER-NORMAL</codigo>
      <valor>25.90</valor>
      <prazo>3</prazo>
      <codigo_rastreio_modelo>ER1234567890BR</codigo_rastreio_modelo>
    </opcao>
    <opcao>
      <tipo>expresso</tipo>
      <codigo>ER-EXPRESS</codigo>
      <valor>45.50</valor>
      <prazo>1</prazo>
      <codigo_rastreio_modelo>EX1234567890BR</codigo_rastreio_modelo>
    </opcao>
  </opcoes>
  <mensagem>Cotação realizada com sucesso</mensagem>
</resultado>
```

### Exemplo de Erro
```xml
<?xml version="1.0" encoding="UTF-8"?>
<resultado>
  <status>erro</status>
  <codigo_erro>ER-005</codigo_erro>
  <mensagem>CEP de destino inválido ou não atendido</mensagem>
</resultado>
```

---

## Transportadora B - LogiFretes

**Características:**
- Formato JSON 
- Resposta rápida (0.5-2 segundos)
- Instabilidade frequente (30% de falhas)
- Autenticação via JWT Bearer Token

### Endpoint
```
POST http://localhost:3000/api/cotacoes
```

### Headers
```
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJsb2dpZnJldGVzIiwibmFtZSI6IkxvZ2lGcmV0ZXMgQVBJIiwiaWF0IjoxNjE2MTQ4MDAwfQ.example-token
```

### Payload da Requisição
```json
{
  "origem": {
    "cep": "01310100"
  },
  "destino": {
    "cep": "04538132"
  },
  "pacote": {
    "peso": 1.5,
    "dimensoes": {
      "altura": 20,
      "largura": 30,
      "comprimento": 15
    },
    "valor": 150.00
  },
  "servicos": ["standard", "express", "economic"]
}
```

### Resposta de Sucesso
```json
{
  "status": "success",
  "request_id": "b3f5a7d9-e201-4b7f-91d5-1c6f5d3672ac",
  "cotacoes": [
    {
      "tipo": "standard",
      "codigo": "STD",
      "valor": 28.50,
      "prazo_dias": 3,
      "regiao_disponivel": true
    },
    {
      "tipo": "express",
      "codigo": "EXP",
      "valor": 42.75,
      "prazo_dias": 1,
      "regiao_disponivel": true
    },
    {
      "tipo": "economic",
      "codigo": "ECO",
      "valor": 19.99,
      "prazo_dias": 5,
      "regiao_disponivel": true
    }
  ],
  "meta": {
    "processado_em": "0.1523s"
  }
}
```

### Exemplo de Erro (HTTP 503 Service Unavailable)
```json
{
  "status": "error",
  "error_code": "SERVICE_UNAVAILABLE",
  "message": "Serviço temporariamente indisponível. Tente novamente mais tarde.",
  "request_id": "7ac9f430-d68a-4d7c-9c52-68b1284f512a"
}
```

### Exemplo de Erro (HTTP 400 Bad Request)
```json
{
  "status": "error",
  "error_code": "VALIDATION_ERROR",
  "message": "Erro de validação nos dados informados",
  "details": [
    {
      "field": "destino.cep",
      "message": "CEP de destino inválido"
    }
  ],
  "request_id": "e8f1a293-7461-4f3b-8ac3-9f3d7e51d2b9"
}
```

---

## Transportadora C - MegaShipping

**Características:**
- API REST com parâmetros na URL (GET)
- Limite de 5 requisições por minuto
- Autenticação via API Key como parâmetro
- Restrição de atendimento por regiões

### Endpoint
```
GET http://localhost:3000/shipping/quote
```

### Parâmetros de URL
```
?api_key=MS-A12B34C56D78E90F
&origin=01310100
&destination=04538132
&weight=1.5
&length=15
&width=30
&height=20
&declared_value=150.00
&service_type=all
```

### Resposta de Sucesso
```json
{
  "status": "success",
  "quote_id": "MS-789456123",
  "services": [
    {
      "service_code": "MS-STD",
      "service_name": "Standard",
      "price": 32.10,
      "delivery_time": {
        "min_days": 2,
        "max_days": 3,
        "estimated_date": "2025-05-09"
      },
      "restrictions": []
    },
    {
      "service_code": "MS-EXP",
      "service_name": "Express",
      "price": 50.90,
      "delivery_time": {
        "min_days": 1,
        "max_days": 1,
        "estimated_date": "2025-05-07"
      },
      "restrictions": []
    }
  ],
  "available_regions": ["sudeste", "sul", "centro-oeste"],
  "rate_limit": {
    "remaining": 4,
    "reset_in_seconds": 60,
    "daily_quota": 100,
    "daily_remaining": 97
  }
}
```

### Resposta para Região Não Atendida
```json
{
  "status": "region_not_available",
  "quote_id": "MS-789456124",
  "message": "Região de entrega não atendida por nossos serviços",
  "rate_limit": {
    "remaining": 4,
    "reset_in_seconds": 58,
    "daily_quota": 100,
    "daily_remaining": 96
  }
}
```

### Resposta para Limite Excedido (HTTP 429 Too Many Requests)
```json
{
  "status": "error",
  "error_code": "RATE_LIMIT_EXCEEDED",
  "message": "Limite de requisições excedido",
  "rate_limit": {
    "remaining": 0,
    "reset_in_seconds": 37,
    "daily_quota": 100,
    "daily_remaining": 95
  }
}
```