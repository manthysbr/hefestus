# Hefestus API 

[![Go Version](https://img.shields.io/github/go-mod/go-version/yourusername/hefestus)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Hefestus is a Go-powered API that leverages local Language Models (via Ollama) to analyze and resolve development errors across different domains (Kubernetes, GitHub Actions, ArgoCD). It provides smart, context-aware solutions for common development issues.

## 🌟 Features

- **Domain-Specific Error Analysis**: Specialized handling for Kubernetes, GitHub Actions, and ArgoCD errors
- **Smart Error Analysis**: Get concise root cause analysis and detailed solutions
- **Local LLM Integration**: Uses Ollama for fast, private error resolution
- **Pattern Matching**: Uses pre-defined error patterns for better solutions
- **Swagger Documentation**: Interactive API documentation
- **JSON Responses**: Clean, structured response format

## 🛠️ Prerequisites

- Go 1.21+
- [Ollama](https://ollama.ai/) with a compatible model (e.g., qwen, mistral)
- Docker (optional)

## 🚀 Quick Start

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/hefestus.git
cd hefestus
```

2. **Set up environment**
```bash
cp .env.example .env
# Edit .env with your settings
```

3. **Run Ollama**
```bash
ollama run qwen2.5:1.5b  # or your preferred model
```

4. **Start the API**
```bash
go run cmd/server/main.go
```

## 📚 API Usage

### Error Resolution Endpoints

#### Kubernetes Error
```bash
curl -X POST http://localhost:8080/api/errors/kubernetes \
  -H "Content-Type: application/json" \
  -d '{
    "error_details": "0/3 nodes are available: insufficient memory",
    "context": "Deploying new pod in production cluster"
  }'
```

#### GitHub Actions Error
```bash
curl -X POST http://localhost:8080/api/errors/github \
  -H "Content-Type: application/json" \
  -d '{
    "error_details": "Error: permission denied to access repository",
    "context": "GitHub Actions workflow execution"
  }'
```

#### ArgoCD Error
```bash
curl -X POST http://localhost:8080/api/errors/argocd \
  -H "Content-Type: application/json" \
  -d '{
    "error_details": "sync failed: connection refused",
    "context": "ArgoCD application sync"
  }'
```

Response Format:
```json
{
  "error": {
    "causa": "Nodes sem memória disponível",
    "solucao": "kubectl describe nodes\nkubectl top nodes\nkubectl edit deployment/app-name"
  },
  "message": "Resolution retrieved successfully"
}
```

## 📁 Project Structure

```
hefestus/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── models/          # Data structures
│   └── services/        # Business logic
├── pkg/
│   └── ollama/          # LLM client
├── config/
│   └── domains.json     # Domain configurations
├── data/
│   └── patterns/        # Error pattern dictionaries
│       ├── kubernetes_errors.json
│       ├── github_errors.json
│       └── argocd_errors.json
└── api/                 # API client library
```

[Rest of the README remains the same...]
```

## 📔 Swagger

Here are some prints how the swagger should look like:
> [!IMPORTANT]
> ![swagger1](https://github.com/manthysbr/hefestus/blob/main/img/image1.png)


> [!IMPORTANT]
> ![swagger2](https://github.com/manthysbr/hefestus/blob/main/img/image2.png)

## Using it with Docker 🐳

Run the command below:

```
docker build -t hefestus:latest .
```

Then run it locally:
```
docker run -d \
    -p 8080:8080 \
    -e OLLAMA_MODEL=qwen2.5:1.5b \
    --name hefestus \
    hefestus:latest
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Submit a Pull Request ( or just copy it and paste lol )

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details but is open for everyone to use it.
```
