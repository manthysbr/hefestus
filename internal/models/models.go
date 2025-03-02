package models

// APIError representa um erro padronizado da API
// @Description Estrutura de erro padrão retornada pela API
type APIError struct {
	Code    int    `json:"code" example:"400" binding:"required"`
	Message string `json:"message" example:"Parâmetros inválidos" binding:"required"`
	Details string `json:"details,omitempty" example:"O campo error_details é obrigatório"`
}

// ErrorRequest representa uma solicitação para analisar um erro
// @Description Requisição contendo os detalhes do erro a ser analisado
type ErrorRequest struct {
	ErrorDetails string `json:"error_details" validate:"required" example:"CrashLoopBackOff: container failed to start" binding:"required"`
	Context      string `json:"context" example:"Deployment em cluster Kubernetes 1.26 com imagem Docker personalizada"`
}

// ErrorResponse representa a resposta da API com a solução do erro
// @Description Resposta contendo análise e solução para o erro reportado
type ErrorResponse struct {
	Error   *ErrorSolution `json:"error" binding:"required"`
	Message string         `json:"message" example:"Análise concluída com sucesso"`
}

// ErrorSolution contém a causa raiz e a solução do erro
// @Description Estrutura contendo a causa identificada e soluções propostas para o erro
type ErrorSolution struct {
	Causa   string `json:"causa" example:"Imagem Docker inválida" binding:"required"`
	Solucao string `json:"solucao" example:"kubectl describe pod meu-pod\nkubectl logs meu-pod --previous" binding:"required"`
}
type DomainConfig struct {
	Name           string                 `json:"name" example:"GitHub Actions"`
	PromptTemplate string                 `json:"prompt_template"`
	Parameters     map[string]interface{} `json:"parameters"`
	DictionaryPath string                 `json:"dictionary_path"`
}
