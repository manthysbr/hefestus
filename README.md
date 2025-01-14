# Hefestus API 🚀

[![Go Version](https://img.shields.io/github/go-mod/go-version/yourusername/hefestus)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](#)
[![Swagger Docs](https://img.shields.io/badge/docs-Swagger-85EA2D)](#📖-documentação)

**Hefestus** é uma API escrita em Go que utiliza modelos de linguagem (LLMs) locais via **Ollama** para analisar e resolver erros de desenvolvimento em múltiplos domínios, como **Kubernetes**, **GitHub Actions**, e **ArgoCD**.

---

## 🎯 Objetivo

O Hefestus foi desenvolvido como um projeto de estudo para explorar:
- Como utilizar **Golang** para criar uma API de troubleshooting de erros de desenvolvimento.
- Ideias inovadoras para integrar detecção de erros e automação de solução em pipelines de CI/CD e outras ferramentas (e.g., Teams, Slack, GitHub).

Você pode configurar o Hefestus para receber erros de endpoints ou pipelines e obter soluções diretamente no console ou em outros sistemas integrados.

---

## 🌟 Principais Recursos

| Funcionalidade               | Descrição                                                                                       |
|------------------------------|-------------------------------------------------------------------------------------------------|
| **Arquitetura Multi-Domínio**| Processamento dinâmico com prompts especializados para cada domínio (e.g., Kubernetes, GitHub). |
| **Integração com LLMs**      | Usa **Ollama** para processar modelos open-source localmente com prompts otimizados.            |
| **Dicionários de Erros**     | Banco de dados com padrões de erros por domínio, aumentando a precisão das soluções.           |
| **Swagger UI**               | Documentação interativa para testar endpoints da API.                                          |
| **Controle de Parâmetros**   | Ajustes finos por domínio: temperatura, tokens, etc.                                           |

---

## 🛠️ Tecnologias Utilizadas

- **Go 1.21+**: Linguagem de programação principal.
- **Gin**: Framework web para construção da API.
- **Ollama**: Para processar modelos de linguagem localmente.
- **Swagger**: Documentação interativa da API.
- **Docker**: Containerização para fácil deploy.

---

## 📚 Estrutura do Projeto

```plaintext
hefestus/
├── cmd/server/          # Ponto de entrada da aplicação
├── internal/
│   ├── models/         # Estruturas de dados
│   └── services/       # Lógica de negócio
├── pkg/
│   └── ollama/         # Cliente LLM
├── config/
│   └── domains.json    # Configurações de domínio
└── data/patterns/      # Dicionários de erros
```

---

## 🚀 Como Usar?

### Pré-requisitos

Certifique-se de ter as seguintes dependências instaladas:
- **Go** (1.21+)
- **Docker** (para rodar a aplicação em contêiner, opcional)

### Instalação

```bash
git clone https://github.com/yourusername/hefestus.git
cd hefestus
cp .env.example .env
go run cmd/server/main.go
```

---

## 🔍 Exemplos de Uso

### **Exemplo de Erro no Kubernetes**
```bash
curl -X POST http://localhost:8080/api/errors/kubernetes \
  -H "Content-Type: application/json" \
  -d '{
    "error_details": "0/3 nodes are available: insufficient memory",
    "context": "Deploying new pod in production cluster"
  }'
```

#### **Resposta**
```json
{
  "error": {
    "causa": "Nodes sem memória disponível",
    "solucao": "kubectl describe nodes\nkubectl top nodes"
  },
  "message": "Resolution retrieved successfully"
}
```

### **Swagger UI**
Acesse a documentação interativa no navegador:
```
http://localhost:8080/swagger/index.html
```

---

## 🐳 Docker

### Build da Imagem
```bash
docker build -t hefestus:latest .
```

### Executar o Container
```bash
docker run -d \
    -p 8080:8080 \
    -e OLLAMA_MODEL=qwen2.5:1.5b \
    --name hefestus \
    hefestus:latest
```

---

## 🔧 Configurações Internas

### Estrutura de Chaveamento por Domínio

Exemplo de configuração de domínio (`domains.json`):
```json
{
  "domains": [
    {
      "name": "kubernetes",
      "prompt_template": "Analyze the Kubernetes error and suggest solutions.",
      "parameters": {
        "temperature": 0.7,
        "max_tokens": 150
      },
      "dictionary_path": "data/patterns/kubernetes.json"
    }
  ]
}
```

### Exemplo de Padrão de Erro
```json
{
  "patterns": {
    "insufficient_resources": {
      "pattern": "\\b(insufficient|not enough)\\s+(cpu|memory|resources)\\b",
      "category": "RESOURCE_LIMITS",
      "solutions": [
        "Verifique o uso de recursos no cluster.",
        "Considere aumentar os recursos alocados."
      ]
    }
  }
}
```

---

## 💡 Contribuindo

Contribuições são bem-vindas! Para contribuir:
1. Faça um fork do repositório.
2. Crie um branch para sua feature ou correção: `git checkout -b minha-feature`.
3. Faça um commit das suas alterações: `git commit -m 'Adiciona minha nova feature'`.
4. Faça o push para o branch: `git push origin minha-feature`.
5. Abra um pull request.

---

## 📝 Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

---

## 📫 Contato

Se você tiver dúvidas, entre em contato via [GitHub Issues](https://github.com/yourusername/hefestus/issues) ou envie um e-mail para `email@domain.com`.

---