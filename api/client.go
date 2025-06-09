// Package client fornece um cliente HTTP para interagir com a API Hefestus.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hefestus-api/internal/models"
)

// Client encapsula um cliente HTTP para comunicação com a API Hefestus.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// ClientOption define uma função para configurar opções do cliente.
type ClientOption func(*Client)

// WithTimeout configura o timeout do cliente HTTP.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithBaseURL configura a URL base para o cliente.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// New cria uma nova instância do cliente Hefestus com as opções fornecidas.
func New(options ...ClientOption) *Client {
	client := &Client{
		baseURL: "http://localhost:8080/api",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Aplica as opções de configuração
	for _, option := range options {
		option(client)
	}

	return client
}

// SendErrorRequest envia um pedido de análise de erro para a API e retorna as resoluções sugeridas.
func (c *Client) SendErrorRequest(ctx context.Context, request models.ErrorRequest) (*models.ErrorResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("falha ao serializar requisição: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/errors", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status inesperado: %d, resposta: %s", resp.StatusCode, body)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	return &response, nil
}

// HealthCheck verifica o status de saúde da API.
func (c *Client) HealthCheck(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return "", fmt.Errorf("falha ao criar requisição: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("falha na requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status inesperado: %d, resposta: %s", resp.StatusCode, body)
	}

	var status struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return "", fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	return status.Status, nil
}
