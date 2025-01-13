package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (c *Client) Query(ctx context.Context, errorDetails string, domain string, errorContext string) (string, string, error) {
	// Domain specific prompts
	prompts := map[string]string{
		"kubernetes": "Você é um especialista em Kubernetes e Cloud Native.",
		"github":     "Você é um especialista em CI/CD e GitHub Actions.",
		"argocd":     "Você é um especialista em GitOps e ArgoCD.",
	}

	expertPrompt := prompts[domain]
	if expertPrompt == "" {
		expertPrompt = "Você é um especialista em resolução de erros."
	}

	prompt := fmt.Sprintf(`%s Analise o erro e responda EXATAMENTE neste formato:

CAUSA: [causa em 15 palavras]
SOLUCAO: [solução detalhada]

ERRO: %s
CONTEXTO: %s

Use exatamente CAUSA: e SOLUCAO: como separadores.`,
		expertPrompt,
		errorDetails,
		errorContext)

	// Add logging
	log.Printf("Sending prompt to LLM: %s", prompt)

	reqBody := Request{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.4,
			"top_k":       1.0,
			"top_p":       1.1,
			"max_tokens":  256,
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

	// Log raw response for debugging
	log.Printf("Raw LLM response: %s", apiResponse.Response)

	if apiResponse.Error != "" {
		return "", "", fmt.Errorf("LLM error: %s", apiResponse.Error)
	}

	// Better response parsing with more detailed error handling
	causaStart := strings.Index(apiResponse.Response, "CAUSA:")
	solucaoStart := strings.Index(apiResponse.Response, "SOLUCAO:")

	if causaStart == -1 || solucaoStart == -1 {
		log.Printf("Invalid format detected. Response: %s", apiResponse.Response)
		return "", "", fmt.Errorf("invalid response format from LLM")
	}

	// Extract parts with better bounds checking
	causa := ""
	solucao := ""

	if solucaoStart > causaStart {
		causa = strings.TrimSpace(apiResponse.Response[causaStart+6 : solucaoStart])
	} else {
		return "", "", fmt.Errorf("invalid response format: SOLUCAO appears before CAUSA")
	}

	solucao = strings.TrimSpace(apiResponse.Response[solucaoStart+8:])

	// Validate extracted parts
	if causa == "" || solucao == "" {
		log.Printf("Empty parts detected. Causa: '%s', Solucao: '%s'", causa, solucao)
		return "", "", fmt.Errorf("empty causa or solucao from LLM")
	}

	// Enforce causa length limit
	words := strings.Fields(causa)
	if len(words) > 10 {
		causa = strings.Join(words[:10], " ")
	}

	return causa, solucao, nil
}
