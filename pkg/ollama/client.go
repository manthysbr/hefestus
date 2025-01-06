package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Client struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

type Request struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type Response struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
	Done     bool   `json:"done"`
}

func NewClient() *Client {
	return &Client{
		baseURL:    "http://localhost:11434",
		model:      os.Getenv("OLLAMA_MODEL"),
		httpClient: &http.Client{},
	}
}

func (c *Client) Query(ctx context.Context, errorDetails string, context string) (string, string, error) {
	prompt := fmt.Sprintf(`Analise o erro e forneça uma resposta no formato abaixo:

CAUSA: [MÁXIMO 10 PALAVRAS identificando a causa raiz do erro]
SOLUCAO: [Explicação detalhada incluindo: diagnóstico completo e passos para resolver]

ERRO: %s
CONTEXTO: %s

REGRAS IMPORTANTES:
- CAUSA deve ter NO MÁXIMO 10 palavras, ser extremamente direta
- CAUSA deve identificar apenas o problema raiz
- SOLUCAO deve conter toda explicação detalhada
- Sempre em pt-br
- Manter exatamente o formato CAUSA: e SOLUCAO:
- Não incluir textos conversacionais ou introdutórios`,
		errorDetails,
		context)

	reqBody := Request{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.3, // Reduced for more consistent output
			"top_k":       10,  // Reduced for more focused responses
			"top_p":       0.9,
			"max_tokens":  500,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse Response
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResponse.Error != "" {
		return "", "", fmt.Errorf("LLM error: %s", apiResponse.Error)
	}

	causaStart := strings.Index(apiResponse.Response, "CAUSA:")
	solucaoStart := strings.Index(apiResponse.Response, "SOLUCAO:")

	if causaStart == -1 || solucaoStart == -1 {
		return "", "", fmt.Errorf("invalid response format from LLM")
	}

	causa := strings.TrimSpace(apiResponse.Response[causaStart+6 : solucaoStart])
	solucao := strings.TrimSpace(apiResponse.Response[solucaoStart+8:])

	// Enforce causa length limit
	words := strings.Fields(causa)
	if len(words) > 10 {
		causa = strings.Join(words[:10], " ")
	}

	return causa, solucao, nil
}
