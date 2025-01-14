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

	// Format instructions
	const formatInstructions = `
IMPORTANTE: Responda EXATAMENTE neste formato:

CAUSA: [MÁXIMO 5 PALAVRAS, NUNCA MAIS DO QUE ISSO]
SOLUCAO: [uma sugestão de solução por linha, MÁXIMO ATÉ DUAS SUGESTÕES]

REGRAS:
1. CAUSA deve ter NO MÁXIMO 4 palavras
2. SOLUCAO deve ter comandos ou explicações simples, um por linha
3. Não use && ou comandos compostos
4. Não adicione explicações
5. Sem formatação ou numeração`

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

	// Enhanced response parsing
	response := strings.TrimSpace(apiResponse.Response)
	lines := strings.Split(response, "\n")

	var causa string
	var solucaoLines []string
	var parsingSolucao bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "CAUSA:") {
			causa = strings.TrimSpace(strings.TrimPrefix(line, "CAUSA:"))
			// Validate causa length
			if len(strings.Fields(causa)) > 5 {
				return "", "", fmt.Errorf("causa must have maximum 5 words")
			}
		} else if strings.HasPrefix(line, "SOLUCAO:") {
			parsingSolucao = true
		} else if parsingSolucao {
			// Clean up the command
			line = strings.TrimLeft(line, "0123456789.- *")
			line = strings.TrimSpace(line)
			if line != "" {
				// Split compound commands
				if strings.Contains(line, "&&") {
					commands := strings.Split(line, "&&")
					for _, cmd := range commands {
						cmd = strings.TrimSpace(cmd)
						if cmd != "" {
							solucaoLines = append(solucaoLines, cmd)
						}
					}
				} else {
					solucaoLines = append(solucaoLines, line)
				}
			}
		}
	}

	if causa == "" || len(solucaoLines) == 0 {
		log.Printf("Invalid format detected in response:\n%s", response)
		return "", "", fmt.Errorf("invalid response format from LLM")
	}

	solucao := strings.Join(solucaoLines, "\n")

	return causa, solucao, nil
}
