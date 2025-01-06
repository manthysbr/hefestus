package models

type ErrorRequest struct {
	ErrorDetails string `json:"error_details" example:"go: cannot find module providing package"`
	Context      string `json:"context" example:"Trying to run go build in new project"`
}

type ErrorResponse struct {
	Error   *ErrorSolution `json:"error,omitempty"`
	Message string         `json:"message"`
}

type ErrorSolution struct {
	Causa   string `json:"causa" example:"Módulo Go não inicializado"`
	Solucao string `json:"solucao" example:"Execute go mod init nome-do-projeto"`
}
