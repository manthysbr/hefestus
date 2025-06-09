package integration_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"hefestus-api/internal/models"
)

// Integration tests for Swagger UI and API models
// These tests validate that the refactored models work correctly with the Swagger UI

func TestMain(m *testing.M) {
	// Check if server is running, skip if not
	resp, err := http.Get("http://localhost:8080/api/health")
	if err != nil {
		fmt.Println("Server not running, skipping integration tests")
		fmt.Println("To run these tests, start the server with: go run cmd/server/main.go")
		os.Exit(0)
	}
	resp.Body.Close()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestSwaggerUI_Accessibility(t *testing.T) {
	// Test that Swagger UI is accessible
	resp, err := http.Get("http://localhost:8080/swagger/index.html")
	if err != nil {
		t.Fatalf("Failed to access Swagger UI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected HTML content type, got %s", contentType)
	}
}

func TestSwaggerJSON_Structure(t *testing.T) {
	// Test that Swagger JSON contains our refactored models
	resp, err := http.Get("http://localhost:8080/swagger/doc.json")
	if err != nil {
		t.Fatalf("Failed to access Swagger JSON: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var swaggerDoc map[string]interface{}
	if err := json.Unmarshal(body, &swaggerDoc); err != nil {
		t.Fatalf("Failed to parse Swagger JSON: %v", err)
	}

	// Check that definitions exist
	definitions, ok := swaggerDoc["definitions"].(map[string]interface{})
	if !ok {
		t.Fatal("Definitions not found in Swagger JSON")
	}

	// Verify our models are defined correctly
	expectedModels := []string{
		"models.APIError",
		"models.ErrorRequest",
		"models.ErrorResponse",
		"models.ErrorSolution",
	}

	for _, model := range expectedModels {
		if _, exists := definitions[model]; !exists {
			t.Errorf("Model %s not found in Swagger definitions", model)
		}
	}

	// Check ErrorSolution structure to ensure JSON tags are correct
	errorSolution, ok := definitions["models.ErrorSolution"].(map[string]interface{})
	if !ok {
		t.Fatal("ErrorSolution definition not found")
	}

	properties, ok := errorSolution["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("ErrorSolution properties not found")
	}

	// Verify JSON field names are maintained for backward compatibility
	if _, exists := properties["causa"]; !exists {
		t.Error("Expected 'causa' field in ErrorSolution JSON schema")
	}

	if _, exists := properties["solucao"]; !exists {
		t.Error("Expected 'solucao' field in ErrorSolution JSON schema")
	}
}

func TestHealthCheck_EndpointResponse(t *testing.T) {
	// Test health check endpoint
	resp, err := http.Get("http://localhost:8080/api/health")
	if err != nil {
		t.Fatalf("Failed to call health check: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode health check response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestAPIError_JSONCompatibility(t *testing.T) {
	// Test that APIError serializes correctly for API responses
	apiError := models.NewAPIError(400, "Bad Request", "Invalid parameters")

	jsonData, err := json.Marshal(apiError)
	if err != nil {
		t.Fatalf("Failed to marshal APIError: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal APIError JSON: %v", err)
	}

	// Verify required fields
	if decoded["code"].(float64) != 400 {
		t.Errorf("Expected code 400, got %v", decoded["code"])
	}

	if decoded["message"].(string) != "Bad Request" {
		t.Errorf("Expected message 'Bad Request', got %v", decoded["message"])
	}

	if decoded["details"].(string) != "Invalid parameters" {
		t.Errorf("Expected details 'Invalid parameters', got %v", decoded["details"])
	}
}

func TestErrorSolution_BackwardCompatibility(t *testing.T) {
	// Test that ErrorSolution maintains backward compatibility with JSON field names
	solution := models.NewErrorSolution("Configuration error", "Check your settings")

	jsonData, err := json.Marshal(solution)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorSolution: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ErrorSolution JSON: %v", err)
	}

	// Verify JSON field names are in Portuguese for backward compatibility
	if decoded["causa"].(string) != "Configuration error" {
		t.Errorf("Expected causa 'Configuration error', got %v", decoded["causa"])
	}

	if decoded["solucao"].(string) != "Check your settings" {
		t.Errorf("Expected solucao 'Check your settings', got %v", decoded["solucao"])
	}

	// Verify Go field names (Cause/Solution) are not in JSON
	if _, exists := decoded["cause"]; exists {
		t.Error("Go field name 'cause' should not appear in JSON")
	}

	if _, exists := decoded["solution"]; exists {
		t.Error("Go field name 'solution' should not appear in JSON")
	}
}

func TestErrorResponse_StructureIntegrity(t *testing.T) {
	// Test complete ErrorResponse structure
	solution := models.NewErrorSolution("Pod failure", "Check pod logs")
	response := models.NewErrorResponse(solution, "Analysis completed successfully")

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorResponse: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ErrorResponse JSON: %v", err)
	}

	// Check top-level structure
	if decoded["message"].(string) != "Analysis completed successfully" {
		t.Errorf("Expected message 'Analysis completed successfully', got %v", decoded["message"])
	}

	// Check nested error structure
	errorObj, ok := decoded["error"].(map[string]interface{})
	if !ok {
		t.Fatal("Error object not found in response")
	}

	if errorObj["causa"].(string) != "Pod failure" {
		t.Errorf("Expected causa 'Pod failure', got %v", errorObj["causa"])
	}

	if errorObj["solucao"].(string) != "Check pod logs" {
		t.Errorf("Expected solucao 'Check pod logs', got %v", errorObj["solucao"])
	}
}

func TestDomainConfig_ConstructorFunctionality(t *testing.T) {
	// Test DomainConfig constructor with parameters
	params := map[string]interface{}{
		"temperature": 0.7,
		"max_tokens":  150,
	}

	config := models.NewDomainConfig("Kubernetes", "Analyze k8s errors", "/path/to/dict", params)

	if config.Name != "Kubernetes" {
		t.Errorf("Expected name 'Kubernetes', got %s", config.Name)
	}

	if config.PromptTemplate != "Analyze k8s errors" {
		t.Errorf("Expected prompt template 'Analyze k8s errors', got %s", config.PromptTemplate)
	}

	if config.DictionaryPath != "/path/to/dict" {
		t.Errorf("Expected dictionary path '/path/to/dict', got %s", config.DictionaryPath)
	}

	if config.Parameters["temperature"].(float64) != 0.7 {
		t.Errorf("Expected temperature 0.7, got %v", config.Parameters["temperature"])
	}

	// Test constructor with nil parameters
	configNil := models.NewDomainConfig("Simple", "Simple template", "/simple", nil)
	if configNil.Parameters == nil {
		t.Error("Expected non-nil parameters map")
	}
	if len(configNil.Parameters) != 0 {
		t.Errorf("Expected empty parameters map, got %v", configNil.Parameters)
	}
}
