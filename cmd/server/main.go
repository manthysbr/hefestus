package main

import (
	"log"
	"os"

	_ "hefestus-api/docs"
	"hefestus-api/internal/handlers"
	"hefestus-api/internal/services"
	"hefestus-api/pkg/ollama"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Hefestus API
// @version 1.0
// @description API para resolução de erros técnicos utilizando LLMs locais
// @host localhost:8080
// @BasePath /api
// @schemes http

func init() {
	// Carrega variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente")
	}
}

func main() {
	// Configura o modo do Gin baseado no ambiente
	setGinMode()

	// Inicializa o router
	r := gin.Default()

	// Inicializa serviços
	dictService, err := services.NewDictionaryService()
	if err != nil {
		log.Fatal("Falha ao inicializar serviço de dicionário:", err)
	}

	// Inicializa cliente Ollama
	ollamaClient := ollama.NewClient()

	// Inicializa serviços
	llmService := services.NewLLMService(ollamaClient, dictService)

	// Inicializa handlers
	errorHandler := handlers.NewErrorHandler(llmService)

	// Configura documentação Swagger
	ConfigureSwagger(r)

	// Configura rotas da API
	api := r.Group("/api")
	{
		api.GET("/health", errorHandler.HealthCheck)
		api.POST("/errors/:domain", errorHandler.AnalyzeError)
	}

	// Inicia o servidor
	port := getPort()
	log.Printf("Servidor iniciado na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Falha ao iniciar servidor: ", err)
	}
}

// setGinMode configura o modo do Gin baseado no ambiente
func setGinMode() {
	env := os.Getenv("ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

// getPort retorna a porta configurada ou a padrão
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
