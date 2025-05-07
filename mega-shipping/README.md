# MegaShipping API Mock

Este é um mock da API da transportadora MegaShipping (Transportadora C) para o desafio de Freight Rate Gateway.

## Características

- API REST com parâmetros na URL (GET)
- Limite de 5 requisições por minuto
- Autenticação via API Key como parâmetro
- Restrição de atendimento por regiões

## Executando a API

### Usando Docker

```bash
# Construir a imagem
docker build -t mega-shipping-api .

# Executar o container
docker run -p 3000:3000 mega-shipping-api
```

### Executando localmente

```bash
# Instalar dependências
npm install

# Iniciar servidor
npm start

# Ou para desenvolvimento com hot-reload
npm run dev
```

## Endpoint de Cotação

```
GET http://localhost:3000/shipping/quote
```

### Parâmetros de URL

| Parâmetro      | Descrição                           | Obrigatório |
|----------------|-------------------------------------|-------------|
| api_key        | Chave de API para autenticação      | Sim         |
| origin         | CEP de origem                       | Sim         |
| destination    | CEP de destino                      | Sim         |
| weight         | Peso do pacote em kg                | Sim         |
| length         | Comprimento do pacote em cm         | Sim         |
| width          | Largura do pacote em cm             | Sim         |
| height         | Altura do pacote em cm              | Sim         |
| declared_value | Valor declarado do conteúdo         | Não         |
| service_type   | Tipo de serviço (all, standard, express) | Não    |

### Exemplo de requisição

```
GET http://localhost:3000/shipping/quote?api_key=MS-A12B34C56D78E90F&origin=01310100&destination=04538132&weight=1.5&length=15&width=30&height=20&declared_value=150.00&service_type=all
```

## Limitação de Taxa

A API possui um limite de 5 requisições por minuto e um total de 100 requisições por dia por chave de API.

## Respostas da API

Consulte o arquivo PAYLOADS.md para exemplos de respostas. 