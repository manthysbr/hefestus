// Package models provides data structures for the Hefestus API.
// These models define the request and response formats used throughout
// the error analysis and resolution system.
package models

// APIError represents a standardized API error response.
// This structure provides consistent error information across all API endpoints.
type APIError struct {
	Code    int    `json:"code" example:"400" binding:"required"`
	Message string `json:"message" example:"Invalid parameters" binding:"required"`
	Details string `json:"details,omitempty" example:"The error_details field is required"`
}

// ErrorRequest represents a request to analyze an error.
// Contains the error details and optional context information.
type ErrorRequest struct {
	ErrorDetails string `json:"error_details" validate:"required" example:"CrashLoopBackOff: container failed to start" binding:"required"`
	Context      string `json:"context" example:"Deployment in Kubernetes 1.26 cluster with custom Docker image"`
}

// ErrorResponse represents the API response containing the error solution.
// This is the main response structure returned after error analysis.
type ErrorResponse struct {
	Error   *ErrorSolution `json:"error" binding:"required"`
	Message string         `json:"message" example:"Analysis completed successfully"`
}

// ErrorSolution contains the root cause and solution for an error.
// The field names use Portuguese in JSON tags to maintain API compatibility
// while using English field names following Go conventions.
type ErrorSolution struct {
	Cause    string `json:"causa" example:"Invalid Docker image" binding:"required"`
	Solution string `json:"solucao" example:"kubectl describe pod my-pod\nkubectl logs my-pod --previous" binding:"required"`
}

// DomainConfig represents configuration for a specific technical domain.
// Each domain has its own prompt template, parameters, and error dictionary.
type DomainConfig struct {
	Name           string                 `json:"name" example:"GitHub Actions"`
	PromptTemplate string                 `json:"prompt_template"`
	Parameters     map[string]interface{} `json:"parameters"`
	DictionaryPath string                 `json:"dictionary_path"`
}

// NewAPIError creates a new APIError with the provided parameters.
// This constructor ensures consistent error creation throughout the application.
func NewAPIError(code int, message, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewErrorSolution creates a new ErrorSolution with the provided cause and solution.
// This constructor provides a convenient way to create error solutions.
func NewErrorSolution(cause, solution string) *ErrorSolution {
	return &ErrorSolution{
		Cause:    cause,
		Solution: solution,
	}
}

// NewErrorResponse creates a new ErrorResponse with the provided solution and message.
// This constructor ensures consistent response creation.
func NewErrorResponse(solution *ErrorSolution, message string) *ErrorResponse {
	return &ErrorResponse{
		Error:   solution,
		Message: message,
	}
}

// NewDomainConfig creates a new DomainConfig with the provided parameters.
// This constructor provides a convenient way to create domain configurations.
func NewDomainConfig(name, promptTemplate, dictionaryPath string, parameters map[string]interface{}) *DomainConfig {
	if parameters == nil {
		parameters = make(map[string]interface{})
	}
	return &DomainConfig{
		Name:           name,
		PromptTemplate: promptTemplate,
		Parameters:     parameters,
		DictionaryPath: dictionaryPath,
	}
}
