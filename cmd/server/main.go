package main

// @title Hefestus API
// @version 1.0
// @description Error resolution API using local LLM
// @host localhost:8080
// @BasePath /api
// @schemes http

// @tag.name errors
// @tag.description Error resolution endpoints

import (
	"log"

	_ "hefestus-api/docs"
	"hefestus-api/internal/models"
	"hefestus-api/internal/services"
	"hefestus-api/pkg/ollama"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// @Summary      Analisar e resolver erros por domínio
// @Description  Recebe detalhes de um erro e seu contexto, retornando possíveis soluções baseadas em LLM
// @Tags         errors
// @Accept       json
// @Produce      json
// @Param        domain   path      string                 true  "Domínio técnico (kubernetes, github, argocd)"   Enums(kubernetes, github, argocd)
// @Param        request  body      models.ErrorRequest    true  "Detalhes do erro e contexto"
// @Success      200      {object}  models.ErrorResponse   "Solução para o erro"
// @Failure      400      {object}  models.APIError        "Erro de validação ou requisição inválida"
// @Failure      404      {object}  models.APIError        "Domínio não encontrado"
// @Failure      500      {object}  models.APIError        "Erro interno do servidor"
// @Router       /errors/{domain} [post]
func main() {
	r := gin.Default()

	dictService, err := services.NewDictionaryService()
	if err != nil {
		log.Fatal("Failed to initialize dictionary service:", err)
	}

	ollamaClient := ollama.NewClient()
	llmService := services.NewLLMService(ollamaClient, dictService)

	// Swagger routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Static("/swagger-ui", "./docs/swagger")

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		api.POST("/errors/:domain", func(c *gin.Context) {
			domain := c.Param("domain")
			var request models.ErrorRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			resolution, err := llmService.GetResolution(
				c.Request.Context(),
				domain,
				request.ErrorDetails,
				request.Context,
			)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, models.ErrorResponse{
				Error:   resolution,
				Message: "Resolution retrieved successfully",
			})
		})
	}

	if err := r.Run(); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
