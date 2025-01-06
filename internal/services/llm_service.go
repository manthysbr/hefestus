package services

import (
	"context"
	"fmt"
	"hefestus-api/internal/models"
	"hefestus-api/pkg/ollama"
)

type LLMService struct {
	ollamaClient *ollama.Client
}

func NewLLMService(client *ollama.Client) *LLMService {
	return &LLMService{
		ollamaClient: client,
	}
}

func (s *LLMService) GetResolution(ctx context.Context, errorDetails string, context string) (*models.ErrorSolution, error) {
	causa, solucao, err := s.ollamaClient.Query(ctx, errorDetails, context)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter resolução do LLM: %w", err)
	}

	return &models.ErrorSolution{
		Causa:   causa,
		Solucao: solucao,
	}, nil
}
