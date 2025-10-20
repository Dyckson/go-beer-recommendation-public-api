# 📚 Documentação da API

Esta documentação fornece exemplos detalhados de como usar todos os endpoints da API.

## 🌐 Base URL

```
http://localhost:1112/api
```

## 🍺 Estilos de Cerveja (CRUD)

### 📋 Listar Todos os Estilos

**Endpoint:**
```http
GET /api/beer-styles/list
```

**Exemplo de Requisição:**
```bash
curl -X GET http://localhost:1112/api/beer-styles/list
```

**Resposta de Sucesso (200):**
```json
{
  "beerStyles": [
    {
      "uuid": "123e4567-e89b-12d3-a456-426614174000",
      "name": "IPA",
      "tempMin": -6.0,
      "tempMax": 7.0,
      "createdAt": "2025-10-02T10:00:00Z",
      "updatedAt": "2025-10-02T10:00:00Z"
    },
    {
      "uuid": "987fcdeb-51d3-12a4-a456-426614174001", 
      "name": "Weissbier",
      "tempMin": -1.0,
      "tempMax": 3.0,
      "createdAt": "2025-10-02T10:30:00Z",
      "updatedAt": "2025-10-02T10:30:00Z"
    }
  ]
}
```

**Resposta Vazia (200):**
```json
{
  "beerStyles": []
}
```

**Erro de Servidor (500):**
```json
{
  "message": "internal error"
}
```

### ➕ Criar Novo Estilo

**Endpoint:**
```http
POST /api/beer-styles/create
Content-Type: application/json
```

**Exemplo de Requisição:**
```bash
curl -X POST http://localhost:1112/api/beer-styles/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Imperial Stout",
    "temp_min": -10.0,
    "temp_max": 13.0
  }'
```

**Corpo da Requisição:**
```json
{
  "name": "Imperial Stout",
  "temp_min": -10.0,
  "temp_max": 13.0
}
```

**Resposta de Sucesso (201):**
```json
{
  "data": {
    "uuid": "456e7890-e89b-12d3-a456-426614174002",
    "name": "Imperial Stout", 
    "tempMin": -10.0,
    "tempMax": 13.0,
    "createdAt": "2025-10-02T11:00:00Z",
    "updatedAt": "2025-10-02T11:00:00Z"
  }
}
```

**Validação - Nome Obrigatório (400):**
```json
{
  "message": "name is required"
}
```

**Validação - Temperatura Mínima Obrigatória (400):**
```json
{
  "message": "temp_min is required"
}
```

**Validação - Temperatura Máxima Obrigatória (400):**
```json
{
  "message": "temp_max is required"
}
```

**Validação - Faixa de Temperatura Inválida (400):**
```json
{
  "message": "temperature range invalid: minimum temperature cannot be greater than maximum temperature"
}
```

**Conflito - Nome Já Existe (409):**
```json
{
  "message": "beer style with name 'IPA' already exists"
}
```

**JSON Malformado (400):**
```json
{
  "message": "invalid JSON format"
}
```

### ✏️ Atualizar Estilo Existente

**Endpoint:**
```http
PUT /api/beer-styles/edit/{uuid}
Content-Type: application/json
```

**Exemplo de Requisição:**
```bash
curl -X PUT http://localhost:1112/api/beer-styles/edit/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Double IPA",
    "temp_min": -7.0,
    "temp_max": 8.0
  }'
```

**Corpo da Requisição (campos opcionais):**
```json
{
  "name": "Double IPA",
  "temp_min": -7.0,
  "temp_max": 8.0
}
```

**Resposta de Sucesso (200):**
```json
{
  "data": {
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Double IPA",
    "tempMin": -7.0,
    "tempMax": 8.0,
    "createdAt": "2025-10-02T10:00:00Z",
    "updatedAt": "2025-10-02T11:15:00Z"
  }
}
```

**Estilo Não Encontrado (404):**
```json
{
  "message": "beer style not found"
}
```

### 🗑️ Deletar Estilo

**Endpoint:**
```http
DELETE /api/beer-styles/{uuid}
```

**Exemplo de Requisição:**
```bash
curl -X DELETE http://localhost:1112/api/beer-styles/123e4567-e89b-12d3-a456-426614174000
```

**Resposta de Sucesso (200):**
```json
{
  "message": "Beer style deleted successfully"
}
```

**Estilo Não Encontrado (404):**
```json
{
  "message": "beer style not found"
}
```

## 🎵 Recomendação de Playlist

### 🔍 Obter Recomendação Baseada na Temperatura

**Endpoint:**
```http
POST /api/recommendations/suggest
Content-Type: application/json
```

**Exemplo de Requisição:**
```bash
curl -X POST http://localhost:1112/api/recommendations/suggest \
  -H "Content-Type: application/json" \
  -d '{
    "temperature": -7.0
  }'
```

**Corpo da Requisição:**
```json
{
  "temperature": -7.0
}
```

**Resposta de Sucesso (200):**
```json
{
  "beerStyle": "IPA",
  "playlist": {
    "name": "Rock Playlist for IPA",
    "tracks": [
      {
        "name": "Bohemian Rhapsody",
        "artist": "Queen",
        "link": "https://open.spotify.com/track/4u7EnebtmKWzUH433cf5Qv"
      },
      {
        "name": "Stairway to Heaven", 
        "artist": "Led Zeppelin",
        "link": "https://open.spotify.com/track/BQNHGiwUeAXSVb7JC5SqAA"
      },
      {
        "name": "Hotel California",
        "artist": "Eagles", 
        "link": "https://open.spotify.com/track/40riOy7x9W7GXjyGp4pjAv"
      }
    ]
  }
}
```

**Validação - Temperatura Inválida (400):**
```json
{
  "message": "temperature must be between -50 and 50 degrees Celsius"
}
```

**Validação - JSON Malformado (400):**
```json
{
  "message": "invalid request body format"
}
```

**Nenhuma Playlist Encontrada (404):**
```json
{
  "message": "no playlist found for beer style 'Obscure Style'"
}
```

**Erro na Determinação do Estilo (500):**
```json
{
  "message": "Unable to determine suitable beer style"
}
```

**Spotify Indisponível (503):**
```json
{
  "message": "Spotify service is temporarily unavailable"
}
```

**Erro Interno (500):**
```json
{
  "message": "Internal server error"
}
```

## 📊 Exemplos de Fluxo Completo

### Cenário 1: Criando e Testando um Novo Estilo

```bash
# 1. Criar um novo estilo
curl -X POST http://localhost:1112/api/beer-styles/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Pilsner",
    "temp_min": -2.0,
    "temp_max": 4.0
  }'

# 2. Listar todos os estilos para confirmar
curl -X GET http://localhost:1112/api/beer-styles/list

# 3. Testar recomendação com temperatura ideal para Pilsner
curl -X POST http://localhost:1112/api/recommendations/suggest \
  -H "Content-Type: application/json" \
  -d '{"temperature": 1.0}'
```

### Cenário 2: Atualizando um Estilo Existente

```bash
# 1. Listar estilos para pegar um UUID
curl -X GET http://localhost:1112/api/beer-styles/list

# 2. Atualizar o estilo (use um UUID real da resposta anterior)
curl -X PUT http://localhost:1112/api/beer-styles/edit/UUID_AQUI \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Czech Pilsner",
    "temp_max": 5.0
  }'

# 3. Testar recomendação com nova faixa de temperatura
curl -X POST http://localhost:1112/api/recommendations/suggest \
  -H "Content-Type: application/json" \
  -d '{"temperature": 2.5}'
```

### Cenário 3: Testando Validações

```bash
# Teste 1: Nome duplicado
curl -X POST http://localhost:1112/api/beer-styles/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "IPA",
    "temp_min": -5.0,
    "temp_max": 5.0
  }'

# Teste 2: Temperatura inválida
curl -X POST http://localhost:1112/api/beer-styles/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid Beer",
    "temp_min": 10.0,
    "temp_max": 5.0
  }'

# Teste 3: Temperatura extrema na recomendação
curl -X POST http://localhost:1112/api/recommendations/suggest \
  -H "Content-Type: application/json" \
  -d '{"temperature": 100.0}'
```

## 🔧 Códigos de Status Detalhados

| Código | Significado | Quando Ocorre |
|--------|-------------|---------------|
| **200** | OK | Operação realizada com sucesso |
| **201** | Created | Recurso criado com sucesso |
| **400** | Bad Request | Dados inválidos ou malformados |
| **404** | Not Found | Recurso não encontrado |
| **409** | Conflict | Conflito (ex: nome duplicado) |
| **500** | Internal Server Error | Erro interno do servidor |
| **503** | Service Unavailable | Serviço externo indisponível |

## 🧪 Testando com Diferentes Ferramentas

### cURL (Exemplos acima)

### HTTPie
```bash
# Instalar: pip install httpie

# Listar estilos
http GET localhost:1112/api/beer-styles/list

# Criar estilo
http POST localhost:1112/api/beer-styles/create \
  name="New Style" temp_min:=-5.0 temp_max:=8.0

# Recomendação
http POST localhost:1112/api/recommendations/suggest \
  temperature:=-7.0
```

### Postman
1. Importe a collection de exemplos
2. Configure a variável `base_url` como `http://localhost:1112`
3. Execute os requests na ordem dos cenários

### JavaScript/Fetch
```javascript
// Listar estilos
fetch('http://localhost:1112/api/beer-styles/list')
  .then(response => response.json())
  .then(data => console.log(data));

// Criar estilo
fetch('http://localhost:1112/api/beer-styles/create', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: 'Stout',
    temp_min: -8.0,
    temp_max: 2.0
  })
})
.then(response => response.json())
.then(data => console.log(data));

// Recomendação
fetch('http://localhost:1112/api/recommendations/suggest', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    temperature: -5.0
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 💡 Dicas de Uso

1. **Sempre valide os dados** antes de enviar requisições
2. **Use UUIDs válidos** ao atualizar ou deletar estilos
3. **Teste diferentes temperaturas** para ver o algoritmo de recomendação
4. **Monitore os logs** do servidor para debug em caso de erro
5. **Configure corretamente** as credenciais do Spotify para funcionamento completo

A API está pronta para ser integrada em frontends, aplicações mobile ou outros serviços backend!
