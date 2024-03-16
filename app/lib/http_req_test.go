package lib

import (
	"api/app/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	baseURL := "http://example.com"
	authKey := "myAuthKey"

	// Test with authKey provided
	clientWithAuth := NewHTTPClient(baseURL, authKey)
	if clientWithAuth.authKey != authKey {
		t.Errorf("NewHTTPClient with authKey = %s: authKey = %s; want %s", authKey, clientWithAuth.authKey, authKey)
	}

	// Test without authKey provided
	clientWithoutAuth := NewHTTPClient(baseURL)
	if clientWithoutAuth.authKey != "" {
		t.Errorf("NewHTTPClient without authKey: authKey = %s; want empty string", clientWithoutAuth.authKey)
	}
}

func TestHTTPClient_setHeaders(t *testing.T) {
	LoadEnvironment(config.Environment)

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	client := NewHTTPClient("http://example.com")

	// Test with authKey
	client.authKey = "myAuthKey"
	client.setHeaders(req)
	authHeaderValue := req.Header.Get("Authorization")
	if authHeaderValue != "Bearer myAuthKey" {
		t.Errorf("Authorization header value incorrect, got %s", authHeaderValue)
	}

	// Test without authKey
	client.authKey = ""
	req.Header.Del("Authorization")
	client.setHeaders(req)
	authHeaderValue = req.Header.Get("Authorization")
	if authHeaderValue != "" {
		t.Errorf("Authorization header should be empty, got %s", authHeaderValue)
	}

	// Test setting custom content type
	req.Header.Del("Content-Type")
	client.setHeaders(req, "custom/content-type")
	contentTypeHeaderValue := req.Header.Get("Content-Type")
	if contentTypeHeaderValue != "custom/content-type" {
		t.Errorf("Content-Type header value incorrect, got %s", contentTypeHeaderValue)
	}
	// Add more assertions based on your expected behavior
}

func TestHTTPClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add your test logic here
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	resp, err := client.Get("/some/path")

	if err != nil {
		t.Errorf("Error should be nil, got %v", err)
	}

	if resp == nil {
		t.Error("Response should not be nil")
	}
	// Add more assertions based on your expected response and behavior
}

func TestHTTPClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add your test logic here
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	payload := map[string]interface{}{"key": "value"}
	resp, err := client.Post("/some/path", payload, "application/json")

	if err != nil {
		t.Errorf("Error should be nil, got %v", err)
	}

	if resp == nil {
		t.Error("Response should not be nil")
	}
	// Add more assertions based on your expected response and behavior
}

func TestHTTPClient_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add your test logic here
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	payload := map[string]interface{}{"key": "value"}
	resp, err := client.Put("/some/path", payload, "application/json")

	if err != nil {
		t.Errorf("Error should be nil, got %v", err)
	}

	if resp == nil {
		t.Error("Response should not be nil")
	}
	// Add more assertions based on your expected response and behavior
}

func TestHTTPClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add your test logic here
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	resp, err := client.Delete("/some/path")

	if err != nil {
		t.Errorf("Error should be nil, got %v", err)
	}

	if resp == nil {
		t.Error("Response should not be nil")
	}
	// Add more assertions based on your expected response and behavior
}
