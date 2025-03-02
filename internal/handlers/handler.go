package handlers

import (
	"net/http"

	"hefestus-api/internal/models"
	"hefestus-api/internal/services"

	"github.com/gin-gonic/gin"
)

// ErrorHandler encapsula a manipulação de requisições de análise de erros
type ErrorHandler struct {
	llmService *services.LLMService
}

// NewErrorHandler cria um novo manipulador de erros
func NewErrorHandler(llmService *services.LLMService) *ErrorHandler {
	return &ErrorHandler{
		llmService: llmService,
	}
}

// AnalyzeError processa requisições de análise de erro
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
func (h *ErrorHandler) AnalyzeError(c *gin.Context) {
	domain := c.Param("domain")

	// Validação do domínio
	if !isValidDomain(domain) {
		c.JSON(http.StatusNotFound, models.APIError{
			Code:    http.StatusNotFound,
			Message: "Domínio não encontrado",
			Details: "Domínios válidos: kubernetes, github, argocd",
		})
		return
	}

	var request models.ErrorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    http.StatusBadRequest,
			Message: "Requisição inválida",
			Details: err.Error(),
		})
		return
	}

	// Validar dados da requisição
	if request.ErrorDetails == "" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    http.StatusBadRequest,
			Message: "Campos obrigatórios não preenchidos",
			Details: "O campo error_details é obrigatório",
		})
		return
	}

	// Obter resolução do serviço LLM
	resolution, err := h.llmService.GetResolution(
		c.Request.Context(),
		domain,
		request.ErrorDetails,
		request.Context,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Erro ao processar solução",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ErrorResponse{
		Error:   resolution,
		Message: "Análise concluída com sucesso",
	})
}

// HealthCheck verifica a saúde do serviço
// @Summary      Verificar saúde do serviço
// @Description  Verifica se o serviço está em funcionamento
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (h *ErrorHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// isValidDomain verifica se o domínio solicitado é suportado
func isValidDomain(domain string) bool {
	validDomains := map[string]bool{
		"kubernetes": true,
		"github":     true,
		"argocd":     true,
	}
	return validDomains[domain]
}
