package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
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

type DomainConfig struct {
	Name           string                 `json:"name"`
	PromptTemplate string                 `json:"prompt_template"`
	Parameters     map[string]interface{} `json:"parameters"`
}

func (c *Client) Query(ctx context.Context, errorDetails string, domain string, errorContext string) (string, string, error) {
	// Load domain configuration
	data, err := os.ReadFile("config/domains.json")
	if err != nil {
		return "", "", fmt.Errorf("failed to read domain config: %w", err)
	}

	var config struct {
		Domains map[string]DomainConfig `json:"domains"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return "", "", fmt.Errorf("failed to parse domain config: %w", err)
	}

	domainConfig, exists := config.Domains[domain]
	if !exists {
		return "", "", fmt.Errorf("unknown domain: %s", domain)
	}

	// Add strict format instructions to prompt template
	domainConfig.PromptTemplate = fmt.Sprintf(`%s

REGRAS IMPORTANTES:
1. Responda EXATAMENTE neste formato
2. Use CAUSA: e SOLUCAO: como separadores exatos
3. Não use formatação ou caracteres especiais
4. Separe CAUSA e SOLUCAO com uma única quebra de linha

FORMATO OBRIGATÓRIO:
CAUSA: [máximo 4 palavras]
SOLUCAO: [apenas comandos, um por linha]

ERRO: {{.Error}}
CONTEXTO: {{.Context}}`, domainConfig.PromptTemplate)

	// Parse template
	tmpl, err := template.New("prompt").Parse(domainConfig.PromptTemplate)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	// Execute template
	var promptBuf bytes.Buffer
	err = tmpl.Execute(&promptBuf, map[string]string{
		"Error":   errorDetails,
		"Context": errorContext,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	log.Printf("Sending prompt to LLM: %s", promptBuf.String())

	// Prepare request with stricter parameters
	reqBody := Request{
		Model:  c.model,
		Prompt: promptBuf.String(),
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1,
			"top_p":       0.1,
			"max_tokens":  150,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send request
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

	// Parse response
	var apiResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResponse.Error != "" {
		return "", "", fmt.Errorf("LLM error: %s", apiResponse.Error)
	}

	log.Printf("Raw LLM response: %s", apiResponse.Response)

	// Updated response parsing
	response := strings.TrimSpace(apiResponse.Response)
	parts := strings.Split(response, "\n")

	var causa, solucao string
	var foundCausa, foundSolucao bool

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "CAUSA:") {
			causa = strings.TrimSpace(strings.TrimPrefix(part, "CAUSA:"))
			foundCausa = true
		} else if strings.HasPrefix(part, "SOLUCAO:") {
			solucao = strings.Join(parts[strings.Index(parts, part)+1:], "\n")
			foundSolucao = true
			break
		}
	}

	if !foundCausa || !foundSolucao {
		log.Printf("Invalid format detected. Response: %s", response)
		return "", "", fmt.Errorf("invalid response format from LLM")
	}

	// Clean up any remaining special characters
	causa = strings.Trim(causa, `"* `)
	solucao = strings.Trim(solucao, `"* `)

	return causa, solucao, nil
}
