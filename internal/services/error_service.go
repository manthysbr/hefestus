package services

import (
	"context"
	"errors"
	"hefestus-api/internal/models"
	"log"
)

// ErrorService define o serviço para processamento de erros
type ErrorService struct {
	llmService *LLMService
}

// NewErrorService cria uma nova instância do serviço de erros
func NewErrorService(llmService *LLMService) *ErrorService {
	return &ErrorService{
		llmService: llmService,
	}
}

// ProcessError analisa uma requisição de erro e retorna possíveis soluções
func (s *ErrorService) ProcessError(ctx context.Context, domain string, req models.ErrorRequest) (*models.ErrorResponse, error) {
	// Validação básica
	if req.ErrorDetails == "" {
		return nil, errors.New("error details cannot be empty")
	}

	log.Printf("Processando erro no domínio %s: %s", domain, req.ErrorDetails)

	// Obter resolução através do serviço LLM
	solution, err := s.llmService.GetResolution(ctx, domain, req.ErrorDetails, req.Context)
	if err != nil {
		log.Printf("Erro ao obter resolução: %v", err)
		return nil, err
	}

	return &models.ErrorResponse{
		Error:   solution,
		Message: "Possíveis resoluções recuperadas com sucesso",
	}, nil
}
