package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAPIError_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		apiError APIError
		expected string
	}{
		{
			name: "complete api error",
			apiError: APIError{
				Code:    400,
				Message: "Bad Request",
				Details: "Invalid parameters",
			},
			expected: `{"code":400,"message":"Bad Request","details":"Invalid parameters"}`,
		},
		{
			name: "api error without details",
			apiError: APIError{
				Code:    404,
				Message: "Not Found",
			},
			expected: `{"code":404,"message":"Not Found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.apiError)
			if err != nil {
				t.Fatalf("Failed to marshal APIError: %v", err)
			}

			// Parse both JSON strings to compare as objects to avoid formatting issues
			var expectedObj, actualObj map[string]interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expectedObj); err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}
			if err := json.Unmarshal(jsonData, &actualObj); err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			if !reflect.DeepEqual(expectedObj, actualObj) {
				t.Errorf("Expected %s, got %s", tt.expected, string(jsonData))
			}
		})
	}
}

func TestAPIError_JSONDeserialization(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected APIError
	}{
		{
			name:    "complete api error",
			jsonStr: `{"code":400,"message":"Bad Request","details":"Invalid parameters"}`,
			expected: APIError{
				Code:    400,
				Message: "Bad Request",
				Details: "Invalid parameters",
			},
		},
		{
			name:    "api error without details",
			jsonStr: `{"code":404,"message":"Not Found"}`,
			expected: APIError{
				Code:    404,
				Message: "Not Found",
				Details: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var apiError APIError
			err := json.Unmarshal([]byte(tt.jsonStr), &apiError)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if !reflect.DeepEqual(tt.expected, apiError) {
				t.Errorf("Expected %+v, got %+v", tt.expected, apiError)
			}
		})
	}
}

func TestErrorRequest_JSONSerialization(t *testing.T) {
	tests := []struct {
		name    string
		request ErrorRequest
		wantErr bool
	}{
		{
			name: "complete error request",
			request: ErrorRequest{
				ErrorDetails: "CrashLoopBackOff: container failed to start",
				Context:      "Deployment em cluster Kubernetes 1.26",
			},
			wantErr: false,
		},
		{
			name: "error request without context",
			request: ErrorRequest{
				ErrorDetails: "Error message",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				var unmarshaled ErrorRequest
				if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
					t.Errorf("Failed to unmarshal: %v", err)
				}
				if !reflect.DeepEqual(tt.request, unmarshaled) {
					t.Errorf("Original and unmarshaled don't match: %+v != %+v", tt.request, unmarshaled)
				}
			}
		})
	}
}

func TestErrorSolution_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		solution ErrorSolution
		wantErr  bool
	}{
		{
			name: "complete error solution",
			solution: ErrorSolution{
				Cause:    "Imagem Docker inválida",
				Solution: "kubectl describe pod meu-pod\nkubectl logs meu-pod --previous",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.solution)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				var unmarshaled ErrorSolution
				if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
					t.Errorf("Failed to unmarshal: %v", err)
				}
				if !reflect.DeepEqual(tt.solution, unmarshaled) {
					t.Errorf("Original and unmarshaled don't match: %+v != %+v", tt.solution, unmarshaled)
				}
			}
		})
	}
}

func TestErrorResponse_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		response ErrorResponse
		wantErr  bool
	}{
		{
			name: "complete error response",
			response: ErrorResponse{
				Error: &ErrorSolution{
					Cause:    "Image Pull Error",
					Solution: "Check image name and registry access",
				},
				Message: "Análise concluída com sucesso",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				var unmarshaled ErrorResponse
				if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
					t.Errorf("Failed to unmarshal: %v", err)
				}
				if !reflect.DeepEqual(tt.response, unmarshaled) {
					t.Errorf("Original and unmarshaled don't match: %+v != %+v", tt.response, unmarshaled)
				}
			}
		})
	}
}

func TestDomainConfig_JSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		config DomainConfig
	}{
		{
			name: "complete domain config",
			config: DomainConfig{
				Name:           "GitHub Actions",
				PromptTemplate: "Analyze the GitHub Actions error",
				Parameters: map[string]interface{}{
					"temperature": 0.7,
					"max_tokens":  float64(150), // JSON unmarshals numbers as float64
				},
				DictionaryPath: "data/patterns/github.json",
			},
		},
		{
			name: "minimal domain config",
			config: DomainConfig{
				Name:           "Kubernetes",
				PromptTemplate: "Analyze Kubernetes error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.config)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			var unmarshaled DomainConfig
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("Failed to unmarshal: %v", err)
			}
			if !reflect.DeepEqual(tt.config, unmarshaled) {
				t.Errorf("Original and unmarshaled don't match: %+v != %+v", tt.config, unmarshaled)
			}
		})
	}
}

// Test that required field validation is working for struct tags
func TestErrorRequest_RequiredFields(t *testing.T) {
	// This test validates the struct tags are in place for validation
	request := ErrorRequest{}

	// Get the struct type information
	requestType := reflect.TypeOf(request)

	// Check that ErrorDetails field has required validation tag
	errorDetailsField, exists := requestType.FieldByName("ErrorDetails")
	if !exists {
		t.Error("ErrorDetails field not found")
	}

	validateTag := errorDetailsField.Tag.Get("validate")
	if validateTag != "required" {
		t.Errorf("Expected validate tag 'required', got '%s'", validateTag)
	}

	bindingTag := errorDetailsField.Tag.Get("binding")
	if bindingTag != "required" {
		t.Errorf("Expected binding tag 'required', got '%s'", bindingTag)
	}
}

func TestAPIError_RequiredFields(t *testing.T) {
	// Check that both Code and Message have required binding tags
	apiErrorType := reflect.TypeOf(APIError{})

	codeField, exists := apiErrorType.FieldByName("Code")
	if !exists {
		t.Error("Code field not found")
	}
	if codeField.Tag.Get("binding") != "required" {
		t.Errorf("Expected Code field to have binding:required tag")
	}

	messageField, exists := apiErrorType.FieldByName("Message")
	if !exists {
		t.Error("Message field not found")
	}
	if messageField.Tag.Get("binding") != "required" {
		t.Errorf("Expected Message field to have binding:required tag")
	}
}

// Test constructor functions
func TestNewAPIError(t *testing.T) {
	code := 400
	message := "Bad Request"
	details := "Invalid input"

	apiError := NewAPIError(code, message, details)

	if apiError.Code != code {
		t.Errorf("Expected Code %d, got %d", code, apiError.Code)
	}
	if apiError.Message != message {
		t.Errorf("Expected Message %s, got %s", message, apiError.Message)
	}
	if apiError.Details != details {
		t.Errorf("Expected Details %s, got %s", details, apiError.Details)
	}
}

func TestNewErrorSolution(t *testing.T) {
	cause := "Invalid configuration"
	solution := "Check your config file"

	errorSolution := NewErrorSolution(cause, solution)

	if errorSolution.Cause != cause {
		t.Errorf("Expected Cause %s, got %s", cause, errorSolution.Cause)
	}
	if errorSolution.Solution != solution {
		t.Errorf("Expected Solution %s, got %s", solution, errorSolution.Solution)
	}
}

func TestNewErrorResponse(t *testing.T) {
	solution := &ErrorSolution{
		Cause:    "Test cause",
		Solution: "Test solution",
	}
	message := "Analysis complete"

	response := NewErrorResponse(solution, message)

	if !reflect.DeepEqual(response.Error, solution) {
		t.Errorf("Expected Error %+v, got %+v", solution, response.Error)
	}
	if response.Message != message {
		t.Errorf("Expected Message %s, got %s", message, response.Message)
	}
}

func TestNewDomainConfig(t *testing.T) {
	tests := []struct {
		name           string
		configName     string
		promptTemplate string
		dictionaryPath string
		parameters     map[string]interface{}
	}{
		{
			name:           "with parameters",
			configName:     "Test Domain",
			promptTemplate: "Test template",
			dictionaryPath: "/path/to/dict",
			parameters: map[string]interface{}{
				"temperature": 0.7,
			},
		},
		{
			name:           "without parameters",
			configName:     "Simple Domain",
			promptTemplate: "Simple template",
			dictionaryPath: "/simple/path",
			parameters:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewDomainConfig(tt.configName, tt.promptTemplate, tt.dictionaryPath, tt.parameters)

			if config.Name != tt.configName {
				t.Errorf("Expected Name %s, got %s", tt.configName, config.Name)
			}
			if config.PromptTemplate != tt.promptTemplate {
				t.Errorf("Expected PromptTemplate %s, got %s", tt.promptTemplate, config.PromptTemplate)
			}
			if config.DictionaryPath != tt.dictionaryPath {
				t.Errorf("Expected DictionaryPath %s, got %s", tt.dictionaryPath, config.DictionaryPath)
			}

			// Check parameters
			if tt.parameters == nil {
				if config.Parameters == nil {
					t.Error("Expected empty map, got nil")
				} else if len(config.Parameters) != 0 {
					t.Errorf("Expected empty map, got %v", config.Parameters)
				}
			} else {
				if !reflect.DeepEqual(config.Parameters, tt.parameters) {
					t.Errorf("Expected Parameters %v, got %v", tt.parameters, config.Parameters)
				}
			}
		})
	}
}
