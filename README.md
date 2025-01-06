# Hefestus API

This project is an API built in Go that interacts with a language model running locally on WSL using Ollama. The API is designed to receive error messages from a CI environment and provide possible resolutions using Retrieval-Augmented Generation (RAG) techniques.

## Project Structure

```
error-resolver-api
├── cmd
│   └── server
│       └── main.go          # Entry point of the application
├── internal
│   ├── handlers
│   │   ├── errors_handler.go # Logic for handling error requests
│   │   └── health_handler.go  # Health check logic
│   ├── models
│   │   ├── error_request.go   # Structure for incoming error requests
│   │   └── error_response.go   # Structure for outgoing error responses
│   └── services
│       ├── llm_service.go     # Logic for interacting with the language model
│       └── error_service.go    # Business logic for processing errors
├── api
│   ├── swagger.yaml           # Swagger documentation for the API
│   └── client.go              # Client code for interacting with the API
├── pkg
│   └── ollama
│       └── client.go          # Client code for interacting with the Ollama model
├── go.mod                     # Go module configuration
├── go.sum                     # Checksums for module dependencies
├── Dockerfile                 # Instructions for building a Docker image
├── .env.example               # Example of environment variables
└── README.md                  # Documentation for the project
```

## Setup Instructions

1. **Install Dependencies**: Ensure you have Go, Docker, and WSL installed on your machine.
2. **Clone the Repository**: Clone this repository to your local machine.
3. **Configure Ollama**: Set up Ollama locally in your WSL environment with the necessary models.
4. **Run the Application**: Navigate to the `cmd/server` directory and run `go run main.go` to start the API server.

## Usage

- Use the `/errors` endpoint to send error requests and receive possible resolutions.
- Access the Swagger documentation at `/swagger` to explore the API endpoints and their usage.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.