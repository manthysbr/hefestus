package services

import (
	"context"
	"fmt"
	"hefestus-api/internal/models"
	"hefestus-api/pkg/ollama"
	"strings"
)

type LLMService struct {
	ollamaClient *ollama.Client
	dictService  *DictionaryService
}

func NewLLMService(ollamaClient *ollama.Client, dictService *DictionaryService) *LLMService {
	return &LLMService{
		ollamaClient: ollamaClient,
		dictService:  dictService,
	}
}

func (s *LLMService) GetResolution(ctx context.Context, domain string, errorDetails string, errorContext string) (*models.ErrorSolution, error) {
	// Check dictionary first - Pass both domain and errorDetails
	matches := s.dictService.FindMatches(domain, errorDetails)

	var knownSolutions string
	if len(matches) > 0 {
		knownSolutions = "\nKnown similar errors and solutions:\n"
		for _, match := range matches {
			knownSolutions += fmt.Sprintf("Category: %s\nSolutions: %v\n",
				match.Category, strings.Join(match.Solutions, ", "))
		}
	}

	// Enhanced prompt with dictionary knowledge
	causa, solucao, err := s.ollamaClient.Query(ctx,
		errorDetails+knownSolutions,
		domain,
		errorContext)
	if err != nil {
		return nil, err
	}

	return &models.ErrorSolution{
		Causa:   causa,
		Solucao: solucao,
	}, nil
}
