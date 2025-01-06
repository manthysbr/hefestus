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

// @Summary Get error resolution
// @Description Get possible solutions for an error
// @Tags errors
// @Accept json
// @Produce json
// @Param request body models.ErrorRequest true "Error details and context"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /errors [post]
func main() {
	r := gin.Default()
	ollamaClient := ollama.NewClient()
	llmService := services.NewLLMService(ollamaClient)

	// Swagger routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Static("/swagger-ui", "./docs/swagger")

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		api.POST("/errors", func(c *gin.Context) {
			var request models.ErrorRequest
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			if request.ErrorDetails == "" {
				c.JSON(400, gin.H{"error": "error_details is required"})
				return
			}

			resolution, err := llmService.GetResolution(
				c.Request.Context(),
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
