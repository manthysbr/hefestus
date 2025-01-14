# Hefestus API

[![Go Version](https://img.shields.io/github/go-mod/go-version/yourusername/hefestus)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Hefestus é uma API em Go que utiliza Large Language Models (LLMs) locais através do Ollama para analisar e resolver erros de desenvolvimento em diferentes domínios (Kubernetes, GitHub Actions, ArgoCD).

## 🎯 Objetivo

O Hefestus foi desenvolvido como um projeto de estudo para demonstrar:
- Integração de LLMs locais em aplicações Go
- Arquitetura baseada em domínios
- Processamento contextual de erros
- Documentação automatizada com Swagger
- Padrões de projeto em Go

## 🌟 Características Principais

### Arquitetura Multi-Domínio
- **Chaveamento de Domínios**: Sistema dinâmico que adapta o processamento baseado no domínio (k8s, github, argocd)
- **Prompts Especializados**: Templates específicos para cada domínio
- **Dicionários de Padrões**: Base de conhecimento pré-definida por domínio

### Integração com LLMs e SLMs
- **Ollama**: Processamento local de linguagem natural
- **Prompts Otimizados**: Estruturas de prompt testadas para cada domínio
- **Controle de Parâmetros**: Ajuste fino por domínio (temperatura, tokens, etc)

### Documentação e API
- **Swagger Integrado**: Documentação interativa da API
- **Respostas Padronizadas**: Formato JSON consistente
- **Validação de Entrada**: Verificação de parâmetros e domínios

## 🛠️ Tecnologias Utilizadas

- **Go 1.21+**: Linguagem principal
- **Gin**: Framework web para o swagger da API
- **Ollama**: LLMs locais
- **Swagger**: Documentação da API
- **Docker**: Containerização para a aplicação

## 📚 Estrutura do Projeto

```
hefestus/
├── cmd/server/          # Ponto de entrada da aplicação
├── internal/
│   ├── models/         # Estruturas de dados
│   └── services/       # Lógica de negócio
├── pkg/
│   └── ollama/         # Cliente LLM
├── config/
│   └── 

domains.json

    # Configurações de domínio
└── data/patterns/      # Dicionários de erros
```

## 🔍 Abordagens

### chaveamento de domínio/assunto
```go
type DomainConfig struct {
    Name           string                 `json:"name"`
    PromptTemplate string                 `json:"prompt_template"`
    Parameters     map[string]interface{} `json:"parameters"`
    DictionaryPath string                 `json:"dictionary_path"`
}
```

### error patterns
```json
{
  "patterns": {
    "insufficient_resources": {
      "pattern": "\\b(insufficient|not enough)\\s+(cpu|memory|resources)\\b",
      "category": "RESOURCE_LIMITS",
      "solutions": [...]
    }
  }
}
```

## 🚀 Como eu uso?

### Instalação
```bash
git clone https://github.com/yourusername/hefestus.git
cd hefestus
cp .env.example .env
go run cmd/server/main.go
```

### Exemplos de Uso

**Erro Kubernetes:**
```bash
curl -X POST http://localhost:8080/api/errors/kubernetes \
  -H "Content-Type: application/json" \
  -d '{
    "error_details": "0/3 nodes are available: insufficient memory",
    "context": "Deploying new pod in production cluster"
  }'
```

**Resposta:**
```json
{
  "error": {
    "causa": "Nodes sem memória disponível",
    "solucao": "kubectl describe nodes\nkubectl top nodes"
  },
  "message": "Resolution retrieved successfully"
}
```

## 📖 Documentação

A documentação completa da API está disponível via Swagger UI em:
```
http://localhost:8080/swagger/index.html
```

## 🐳 Docker

Build:
```bash
docker build -t hefestus:latest .
```

Run:
```bash
docker run -d \
    -p 8080:8080 \
    -e OLLAMA_MODEL=qwen2.5:1.5b \
    --name hefestus \
    hefestus:latest
```

## 📝 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo LICENSE para detalhes.
```