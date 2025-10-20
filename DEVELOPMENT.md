# 🛠️ Guia de Desenvolvimento Local

Este guia detalha como configurar e executar o projeto em ambiente de desenvolvimento.

## 🚀 Início Rápido (Docker)

```bash
# 1. Clone e acesse o diretório
git clone <repo-url>
cd backend-test

# 2. Configure as variáveis de ambiente
cp .env.example .env
# Edite o .env com suas credenciais do Spotify

# 3. Execute tudo com Docker
docker-compose up --build

# ✅ API disponível em: http://localhost:1112
```

## 🔧 Desenvolvimento Avançado

### Opção 1: Aplicação Local + Banco Docker

```bash
# 1. Execute apenas o PostgreSQL
docker-compose up database

# 2. Configure variáveis de ambiente
export DATABASE_URL="postgres://beeruser:beerpass@localhost:55432/beerdb?sslmode=disable"
export SPOTIFY_CLIENT_ID="seu_client_id"
export SPOTIFY_CLIENT_SECRET="seu_client_secret"

# 3. Instale dependências
go mod download

# 4. Execute a aplicação
go run main.go

# 5. Execute em modo de desenvolvimento (hot reload)
# Instale air: go install github.com/cosmtrek/air@latest
air
```

### Opção 2: Tudo Local

```bash
# 1. Configure PostgreSQL local
# Instale PostgreSQL e crie:
createdb beerdb
createuser beeruser
# Configure a senha 'beerpass' para o usuário

# 2. Configure variáveis
export DATABASE_URL="postgres://beeruser:beerpass@localhost:5432/beerdb?sslmode=disable"
export SPOTIFY_CLIENT_ID="seu_client_id"
export SPOTIFY_CLIENT_SECRET="seu_client_secret"

# 3. Execute a aplicação
go run main.go
```

## 🧪 Executando Testes

```bash
# Todos os testes
go test ./... -v

# Apenas unit tests
go test ./internal/http/controller/ -v

# Apenas integration tests
go test ./tests/integration/ -v

# Com coverage
go test ./... -cover

# Coverage detalhado
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 🔍 Debugging

### VS Code
1. Instale a extensão Go
2. Use F5 para debug ou:

```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "env": {
                "DATABASE_URL": "postgres://beeruser:beerpass@localhost:55432/beerdb?sslmode=disable",
                "SPOTIFY_CLIENT_ID": "seu_client_id",
                "SPOTIFY_CLIENT_SECRET": "seu_client_secret"
            }
        }
    ]
}
```

### Logs
```bash
# Debug detalhado
GIN_MODE=debug go run main.go

# Logs estruturados
go run main.go 2>&1 | grep "controller="
```

## 🛠️ Ferramentas Úteis

### Hot Reload com Air
```bash
# Instalar
go install github.com/cosmtrek/air@latest

# Executar
air

# Configuração personalizada (.air.toml)
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "tmp/main"
  full_bin = "APP_ENV=dev APP_USER=air ./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["assets", "tmp", "vendor", "tests"]
```

### Análise de Código
```bash
# Linting
go install golang.org/x/lint/golint@latest
golint ./...

# Formatação
go fmt ./...

# Vet (análise estática)
go vet ./...

# Cyclomatic complexity
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
gocyclo -over 10 .
```

## 🐳 Docker para Desenvolvimento

### Rebuild da Imagem
```bash
# Rebuild completo
docker-compose down
docker-compose build --no-cache
docker-compose up

# Rebuild apenas da API
docker-compose build api
docker-compose up api
```

### Acessar Container
```bash
# Shell no container da API
docker-compose exec api /bin/sh

# Shell no container do banco
docker-compose exec database psql -U beeruser -d beerdb
```

### Logs dos Containers
```bash
# Logs da API
docker-compose logs -f api

# Logs do banco
docker-compose logs -f database

# Todos os logs
docker-compose logs -f
```

## 📊 Monitoramento Local

### Health Check Manual
```bash
# Status da API
curl http://localhost:1112/health

# Listar cervejeiras
curl http://localhost:1112/api/beer-styles/list

# Testar recomendação
curl -X POST http://localhost:1112/api/recommendations/suggest \
  -H "Content-Type: application/json" \
  -d '{"temperature": -7}'
```

### Performance Testing
```bash
# Instalar hey
go install github.com/rakyll/hey@latest

# Teste de carga
hey -n 1000 -c 10 http://localhost:1112/api/beer-styles/list
```

## 🔧 Troubleshooting

### Problemas Comuns

#### 1. Erro de Conexão com Banco
```bash
# Verificar se o PostgreSQL está rodando
docker-compose ps database

# Verificar logs do banco
docker-compose logs database

# Resetar o banco
docker-compose down
docker volume rm backend-test_beer_data
docker-compose up database
```

#### 2. Erro de Spotify API
```bash
# Verificar credenciais
echo $SPOTIFY_CLIENT_ID
echo $SPOTIFY_CLIENT_SECRET

# Testar credenciais manualmente
curl -X POST "https://accounts.spotify.com/api/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=$SPOTIFY_CLIENT_ID&client_secret=$SPOTIFY_CLIENT_SECRET"
```

#### 3. Porta em Uso
```bash
# Verificar processo usando a porta 1112
lsof -i :1112

# Matar processo
kill -9 <PID>

# Usar porta diferente
PORT=8080 go run main.go
```

#### 4. Problemas de Dependências
```bash
# Limpar cache do Go
go clean -modcache

# Reinstalar dependências
rm go.sum
go mod download
go mod tidy
```

## 📁 Estrutura para Desenvolvimento

```
backend-test/
├── .vscode/              # Configurações VS Code
├── tmp/                  # Arquivos temporários (air)
├── coverage.out          # Relatório de coverage
├── .env                  # Variáveis locais (não committar)
├── .env.example          # Template de variáveis
└── .air.toml            # Configuração hot reload
```

## 🎯 Próximos Passos

1. **Configure seu editor** com as extensões Go
2. **Execute os testes** para garantir que tudo funciona
3. **Experimente a API** com diferentes dados
4. **Explore o código** seguindo a arquitetura limpa
5. **Adicione novos recursos** seguindo os padrões existentes

---

💡 **Dica**: Use o VS Code com as extensões Go e Docker para uma experiência de desenvolvimento otimizada!
