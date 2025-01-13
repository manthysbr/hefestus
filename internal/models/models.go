package models

type ErrorRequest struct {
	ErrorDetails string `json:"error_details" example:"go: cannot find module providing package"`
	Context      string `json:"context" example:"Trying to run go build in new project"`
}

type ErrorResponse struct {
	Error   *ErrorSolution `json:"error"`
	Message string         `json:"message"`
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
