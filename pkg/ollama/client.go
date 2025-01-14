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

type LLMResponse struct {
	Causa   string   `json:"causa"`
	Solucao []string `json:"solucao"`
}

func isValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func cleanResponse(response string) string {
	// Remove markdown code blocks
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimSuffix(response, "```")
	return strings.TrimSpace(response)
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

	// Format instructions
	const formatInstructions = `
IMPORTANTE: Retorne APENAS um objeto JSON válido no seguinte formato:

{
    "causa": "[máximo 4 palavras]",
    "solucao": ["comando 1", "comando 2"]
}

REGRAS ESTRITAS:
1. Retorne APENAS o JSON, sem markdown ou formatação
2. causa deve ter NO MÁXIMO 4 palavras
3. solucao deve ter array de comandos simples
4. Não use && ou comandos compostos`

	// Combine template with instructions
	promptTemplate := fmt.Sprintf("%s\n%s\n\nERRO: {{.Error}}\nCONTEXTO: {{.Context}}",
		domainConfig.PromptTemplate,
		formatInstructions)

	// Parse template
	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	var promptBuf bytes.Buffer
	err = tmpl.Execute(&promptBuf, map[string]string{
		"Error":   errorDetails,
		"Context": errorContext,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	log.Printf("Sending prompt to LLM: %s", promptBuf.String())

	// Prepare request
	reqBody := Request{
		Model:   c.model,
		Prompt:  promptBuf.String(),
		Stream:  false,
		Options: domainConfig.Parameters,
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

	// Clean and validate response
	cleanedResponse := cleanResponse(apiResponse.Response)
	if !isValidJSON(cleanedResponse) {
		log.Printf("Invalid JSON format detected: %s", cleanedResponse)
		return "", "", fmt.Errorf("invalid JSON response from LLM")
	}

	var llmResponse LLMResponse
	if err := json.Unmarshal([]byte(cleanedResponse), &llmResponse); err != nil {
		log.Printf("Failed to parse LLM response: %v\nCleaned response: %s", err, cleanedResponse)
		return "", "", fmt.Errorf("invalid response format from LLM")
	}

	// Validate response content
	if llmResponse.Causa == "" || len(llmResponse.Solucao) == 0 {
		return "", "", fmt.Errorf("invalid response content from LLM")
	}

	// Validate causa word count
	if len(strings.Fields(llmResponse.Causa)) > 4 {
		return "", "", fmt.Errorf("causa exceeds maximum word count")
	}

	return llmResponse.Causa, strings.Join(llmResponse.Solucao, "\n"), nil
}
