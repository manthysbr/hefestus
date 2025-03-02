package models

type ErrorRequest struct {
	ErrorDetails string `json:"error_details" validate:"required"`
	Context      string `json:"context"`
}

type ErrorResponse struct {
	Resolutions []string `json:"resolutions"`
	Info        string   `json:"info"`
}

type ErrorSolution struct {
	Causa   string `json:"causa" example:"Módulo Go não inicializado"`
	Solucao string `json:"solucao" example:"Execute go mod init nome-do-projeto"`
}

type DomainConfig struct {
	Name           string                 `json:"name" example:"GitHub Actions"`
	PromptTemplate string                 `json:"prompt_template"`
	Parameters     map[string]interface{} `json:"parameters"`
	DictionaryPath string                 `json:"dictionary_path"`
}
