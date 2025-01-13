package services

import (
	"errors"
	"hefestus-api/internal/models"
)

func ProcessError(req models.ErrorRequest) (*models.ErrorResponse, error) {
	if req.ErrorDetails == "" {
		return nil, errors.New("error details cannot be empty")
	}

	solution := &models.ErrorSolution{
		Causa:   "Erro não processado",
		Solucao: "Contate o administrador",
	}

	return &models.ErrorResponse{
		Error:   solution,
		Message: "Possible resolutions retrieved successfully",
	}, nil
}
