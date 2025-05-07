# LogiFretes API Mock

This is a mock implementation of the LogiFretes API for the Freight Rate Gateway challenge. The API simulates the behavior of the LogiFretes shipping quote service, including:

- Fast response time (0.5-2 seconds)
- Instability (30% chance of failure)
- JWT Bearer Token authentication
- JSON request/response format

## Features

- Validates JWT tokens
- Simulates random failures
- Simulates processing time
- Validates request inputs
- Returns shipping quotes in the specified format

## Authentication

The API accepts JWT Bearer tokens with the following format:
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJsb2dpZnJldGVzIiwibmFtZSI6IkxvZ2lGcmV0ZXMgQVBJIiwiaWF0IjoxNjE2MTQ4MDAwfQ.example-token
```

## API Endpoint

```
POST /api/cotacoes
```

### Request Format

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

## Building and Running

### Build the Docker Image

```bash
docker build -t logifretes-api .
```

### Run the Container

```bash
docker run -p 8080:8080 logifretes-api
```

The API will be available at http://localhost:8080/api/cotacoes

## Testing

You can test the API with curl:

```bash
curl -X POST \
  http://localhost:8080/api/cotacoes \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJsb2dpZnJldGVzIiwibmFtZSI6IkxvZ2lGcmV0ZXMgQVBJIiwiaWF0IjoxNjE2MTQ4MDAwfQ.example-token' \
  -d '{
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
}'
``` 