package traefik_filter_on_field

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilterOnField_ServeHTTP(t *testing.T) {
	// Create a mock handler for the next handler in the chain
	mockHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("OK"))
	})

	// Create a new instance of the filter
	config := CreateConfig()
	config.FieldName = "parameter"
	config.ResponseMessage = "Disallowed content"
	config.DisallowedContent = []string{"invalid"}

	filter := &FilterOnField{
		next:   mockHandler,
		config: config,
		name:   "filter",
	}

	// Create a mock request with a valid parameter
	validReq := httptest.NewRequest(http.MethodGet, "/path?parameter=valid", nil)
	validResp := httptest.NewRecorder()
	filter.ServeHTTP(validResp, validReq)

	// Check if the response is as expected
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, validResp.Code)
	}
	if validResp.Body.String() != "OK" {
		t.Errorf("Expected response body %q, got %q", "OK", validResp.Body.String())
	}

	// Create a mock request with an invalid parameter
	invalidReq := httptest.NewRequest(http.MethodGet, "/path?parameter=invalid", nil)
	invalidResp := httptest.NewRecorder()
	filter.ServeHTTP(invalidResp, invalidReq)

	// Check if the response is as expected
	if invalidResp.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, invalidResp.Code)
	}
	if invalidResp.Body.String() != "Disallowed content" {
		t.Errorf("Expected response body %q, got %q", "Disallowed content", invalidResp.Body.String())
	}
}

func TestFilterOnField_Reject_Numeric(t *testing.T) {
	disallowedNumber := "70000001"
	fieldName := "fieldName"
	// Create a mock handler for the next handler in the chain
	mockHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("OK"))
	})

	// Create a new instance of the filter
	config := CreateConfig()
	config.FieldName = fieldName
	config.ResponseMessage = "Disallowed content"
	config.DisallowedContent = []string{disallowedNumber}

	filter := &FilterOnField{
		next:   mockHandler,
		config: config,
		name:   "filter",
	}

	// Create a mock request with a valid parameter
	validReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/path?%s=valid", fieldName), nil)
	validResp := httptest.NewRecorder()
	filter.ServeHTTP(validResp, validReq)

	// Check if the response is as expected
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, validResp.Code)
	}
	if validResp.Body.String() != "OK" {
		t.Errorf("Expected response body %q, got %q", "OK", validResp.Body.String())
	}

	// Create a mock request with an invalid parameter
	invalidReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/path?%s=%s", fieldName, disallowedNumber), nil)
	invalidResp := httptest.NewRecorder()
	filter.ServeHTTP(invalidResp, invalidReq)

	// Check if the response is as expected
	if invalidResp.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, invalidResp.Code)
	}
	if invalidResp.Body.String() != "Disallowed content" {
		t.Errorf("Expected response body %q, got %q", "Disallowed content", invalidResp.Body.String())
	}
}

func TestFilterOnField_Single_Letter_Field(t *testing.T) {
	disallowedValue := "invalit"
	fieldName := "f"
	// Create a mock handler for the next handler in the chain
	mockHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("OK"))
	})

	// Create a new instance of the filter
	config := CreateConfig()
	config.FieldName = fieldName
	config.ResponseMessage = "Disallowed content"
	config.DisallowedContent = []string{disallowedValue}

	filter := &FilterOnField{
		next:   mockHandler,
		config: config,
		name:   "filter",
	}

	// Create a mock request with a valid parameter
	validReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/path?%s=valid", fieldName), nil)
	validResp := httptest.NewRecorder()
	filter.ServeHTTP(validResp, validReq)

	// Check if the response is as expected
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, validResp.Code)
	}
	if validResp.Body.String() != "OK" {
		t.Errorf("Expected response body %q, got %q", "OK", validResp.Body.String())
	}

	// Create a mock request with an invalid parameter
	invalidReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/path?%s=%s", fieldName, disallowedValue), nil)
	invalidResp := httptest.NewRecorder()
	filter.ServeHTTP(invalidResp, invalidReq)

	// Check if the response is as expected
	if invalidResp.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, invalidResp.Code)
	}
	if invalidResp.Body.String() != "Disallowed content" {
		t.Errorf("Expected response body %q, got %q", "Disallowed content", invalidResp.Body.String())
	}
}
