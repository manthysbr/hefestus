package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Hefestus API
// @version         1.0
// @description     API para resolução de erros técnicos utilizando LLMs locais
// @termsOfService  http://swagger.io/terms/

// @contact.name   Suporte Hefestus
// @contact.url    https://github.com/your-username/hefestus/issues
// @contact.email  suporte@hefestus.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api
// @schemes   http

// ConfigureSwagger configura as rotas do Swagger UI
func ConfigureSwagger(r *gin.Engine) {
	url := "/swagger/doc.json"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(url)))
}
