# 🍺 Beer Temperature & Spotify Recommendation API

Uma API REST para gerenciamento de estilos de cerveja e recomendação de playlists baseada na temperatura ideal de consumo.

## Como Executar Localmente

### Pré-requisitos

- Docker e Docker Compose instalados

### Setup Rápido

```bash
# 1. Clone o repositório
git clone https://github.com/your-username/go-beer-recommendation-api.git
cd go-beer-recommendation-api

# 2. 🔒 Configure variáveis de ambiente
cp .env_exemplo .env
# Edite .env com suas credenciais reais do Spotify

# 3. Execute com Docker
docker-compose up --build -d

# ✅ API disponível em: http://localhost:1112
```

## 🔒 Configuração de Segurança

### 🎵 Credenciais do Spotify

1. Acesse: [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Crie uma nova app
3. Copie **Client ID** e **Client Secret**
4. Cole no seu arquivo `.env` (NUNCA no código!)

### 📋 Variáveis de Ambiente

```bash
# Copie o template
cp .env_exemplo .env

# Edite com suas credenciais
SPOTIFY_CLIENT_ID=seu_client_id_aqui
SPOTIFY_CLIENT_SECRET=seu_client_secret_aqui
```

⚠️ **IMPORTANTE**: 
- NUNCA commite arquivos `.env` 
- Use apenas o template `.env_exemplo`
- Mantenha credenciais seguras

## 📚 Como Usar a API

### Base URL

```
http://localhost:1112/api
```

### 🍺 Estilos de Cerveja (CRUD)

#### Listar todos os estilos

```bash
curl -X GET http://localhost:1112/api/beer-styles/list
```

#### Criar novo estilo

```bash
curl -X POST http://localhost:1112/api/beer-styles/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "IPA",
    "temp_min": -6.0,
    "temp_max": 7.0
  }'
```

#### Atualizar estilo

```bash
curl -X PUT http://localhost:1112/api/beer-styles/edit/{uuid} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Double IPA",
    "temp_min": -7.0,
    "temp_max": 8.0
  }'
```

#### Deletar estilo

```bash
curl -X DELETE http://localhost:1112/api/beer-styles/{uuid}
```

### 🎵 Recomendação de Playlist

#### Obter recomendação baseada na temperatura

```bash
curl -X POST http://localhost:1112/api/recommendations/suggest \
  -H "Content-Type: application/json" \
  -d '{"temperature": -7.0}'
```

**Resposta:**

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
      }
    ]
  }
}
```

## 🧪 Executar Testes

```bash
# Todos os testes
go test ./... -v

# Apenas unit tests
go test ./internal/http/controller/ -v

# Apenas integration tests
go test ./tests/integration/ -v
```

## Tecnologias

- **Go 1.24.5** com Gin framework
- **PostgreSQL**
- **Spotify Web API**
- **Docker**
- **Clean Architecture** com testes híbridos

## 📋 Status Codes

| Código | Descrição               |
| ------- | ------------------------- |
| `200` | Sucesso                   |
| `201` | Criado                    |
| `400` | Dados inválidos          |
| `404` | Não encontrado           |
| `409` | Conflito (nome duplicado) |
| `500` | Erro interno              |
| `503` | Spotify indisponível     |

---

**Para documentação completa:** consulte os arquivos `DEVELOPMENT.md`, `API.md` e `FEATURES.md`
