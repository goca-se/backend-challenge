# 🚚 Backend Live Challenge: Freight Rate Gateway

## 🌍 Cenário Caótico

Em um e-commerce de **crescimento explosivo**, o sistema de cotação de fretes se tornou o **pior gargalo** da plataforma. Os clientes abandonam carrinhos em massa porque:

- As cotações demoram **eternos 30 segundos** para carregar
- Transportadoras diferentes mostram **preços inconsistentes**
- O sistema frequentemente **falha completamente** durante picos de tráfego
- Algumas regiões simplesmente **não recebem cotações**

A diretoria ameaça demitir o time inteiro se as taxas de abandono não caírem. **Você é a última esperança** para resolver este caos logístico antes do Black Friday que acontece em 48 horas!

---

## 🧠 Desafio

Construa um **gateway de cotação de fretes** que integre múltiplas transportadoras, garantindo que o cliente sempre receba cotações rápidas e confiáveis.

> Linguagens sugeridas: **Ruby, Node.js, Go**  
> Pode usar bibliotecas open-source (justifique suas escolhas)

---

## 🎯 Requisitos Mínimos

- ✅ **Endpoint para cotação de fretes**
  - Recebe origem, destino, dimensões e peso do pacote
  - Retorna opções de frete disponíveis

- 📊 **Integração com transportadoras**
  - Conectar com 3 transportadoras diferentes
  - Lidar com as particularidades de cada uma

- 🚀 **Performance e disponibilidade**
  - Garantir resposta rápida ao cliente (<3 segundos)
  - Sistema deve funcionar mesmo quando transportadoras falham

- 📈 **Qualidade do resultado**
  - Fornecer cotações precisas e ordenadas
  - Evitar resultados duplicados ou inconsistentes

---

## 🧪 Setup de Teste (Docker Compose)
Para simular essas transportadoras, use o arquivo `docker-compose.yml` que está na raiz do projeto. A partir desse momento, cada transportadora estará disponível e rodando localmente.

> **Transportadoras**:
> - **Transportadora A (EntregasRápidas)**: Responde em XML, alta latência
> - **Transportadora B (LogiFretes)**: API JSON, instabilidade frequente 
> - **Transportadora C (MegaShipping)**: REST com autenticação, limite de requisições baixo

### Executando os serviços

Você pode iniciar todos os serviços das transportadoras com o script de conveniência:

```bash
# Torne o script executável (se necessário)
chmod +x run-services.sh

# Execute o script
./run-services.sh
```

Alternativamente, você pode usar o docker-compose diretamente:

```bash
docker-compose up -d
```

### Endpoints das Transportadoras

Após iniciar os serviços, as transportadoras estarão disponíveis nos seguintes endpoints:

- **Transportadora A (EntregasRápidas)**: http://localhost:6000/api/v1/entregas-rapidas/cotacao
- **Transportadora B (LogiFretes)**: http://localhost:6001/api/cotacoes
- **Transportadora C (MegaShipping)**: http://localhost:6002/shipping/quote

Para mais detalhes sobre os payloads e respostas de cada transportadora, consulte o arquivo [PAYLOADS.md](PAYLOADS.md).

---

## 🔍 Critérios de Avaliação

- **Arquitetura e design**: Organização do código, abstrações criadas
- **Tratamento de problemas**: Como lida com falhas e limitações
- **Performance**: Estratégias para garantir resposta rápida
- **Clareza e manutenibilidade**: Código legível e bem estruturado

---

Sua tarefa é construir um sistema que resolva o problema de cotação de fretes de forma eficiente e escalável. Você tem 1 hora - não precisa implementar uma solução completa, mas deve demonstrar seu pensamento arquitetural e abordagem para resolver os desafios apresentados.