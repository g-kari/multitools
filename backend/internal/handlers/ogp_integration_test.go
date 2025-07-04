package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ogp-verification-service/internal/handlers"
	"ogp-verification-service/internal/models"
)

func TestOGPHandlerIntegration(t *testing.T) {
	handler := handlers.NewOGPHandler()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "Valid URL request",
			requestBody:    models.OGPRequest{URL: "https://example.com"},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.OGPResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if resp.URL != "https://example.com" {
					t.Errorf("Expected URL to be https://example.com, got %s", resp.URL)
				}
				if resp.Timestamp.IsZero() {
					t.Error("Expected timestamp to be set")
				}
			},
		},
		{
			name:           "Empty URL request",
			requestBody:    models.OGPRequest{URL: ""},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				response := string(body)
				if response != "URL is required\n" {
					t.Errorf("Expected 'URL is required\\n', got %q", response)
				}
			},
		},
		{
			name:           "Invalid JSON request",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				response := string(body)
				if response != "Invalid JSON\n" {
					t.Errorf("Expected 'Invalid JSON\\n', got %q", response)
				}
			},
		},
		{
			name:           "Private IP request",
			requestBody:    models.OGPRequest{URL: "http://192.168.1.1"},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, body []byte) {
				expected := "Error fetching OGP data: private IP addresses are not allowed\n"
				response := string(body)
				if response != expected {
					t.Errorf("Expected %q, got %q", expected, response)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/ogp/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.VerifyOGP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}

			tt.checkResponse(t, rr.Body.Bytes())
		})
	}
}

func TestOGPHandlerCORS(t *testing.T) {
	handler := handlers.NewOGPHandler()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/ogp/verify", nil)
	rr := httptest.NewRecorder()

	handler.VerifyOGP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d for OPTIONS request, got %d", http.StatusOK, status)
	}

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}

	for header, expected := range expectedHeaders {
		if value := rr.Header().Get(header); value != expected {
			t.Errorf("Expected header %s to be '%s', got '%s'", header, expected, value)
		}
	}
}

func TestOGPHandlerRateLimiting(t *testing.T) {
	handler := handlers.NewOGPHandler()

	// Make 10 requests (should all succeed)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/ogp/verify", bytes.NewReader([]byte(`{"url":"https://example.com"}`)))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "127.0.0.1:12345" // Same IP for all requests

		rr := httptest.NewRecorder()
		handler.VerifyOGP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Request %d: Expected status %d, got %d", i+1, http.StatusOK, status)
		}
	}

	// 11th request should be rate limited
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ogp/verify", bytes.NewReader([]byte(`{"url":"https://example.com"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"

	rr := httptest.NewRecorder()
	handler.VerifyOGP(rr, req)

	if status := rr.Code; status != http.StatusTooManyRequests {
		t.Errorf("Expected status %d for rate limited request, got %d", http.StatusTooManyRequests, status)
	}

	if body := rr.Body.String(); body != "Rate limit exceeded\n" {
		t.Errorf("Expected 'Rate limit exceeded\\n', got %q", body)
	}
}

func TestOGPHandlerRateLimitingPerIP(t *testing.T) {
	handler := handlers.NewOGPHandler()

	// Make requests from different IPs
	ips := []string{"127.0.0.1:12345", "127.0.0.2:12345", "127.0.0.3:12345"}

	for _, ip := range ips {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/ogp/verify", bytes.NewReader([]byte(`{"url":"https://example.com"}`)))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = ip

		rr := httptest.NewRecorder()
		handler.VerifyOGP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Request from %s: Expected status %d, got %d", ip, http.StatusOK, status)
		}
	}
}

func TestOGPHandlerMethodNotAllowed(t *testing.T) {
	handler := handlers.NewOGPHandler()

	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/v1/ogp/verify", nil)
		rr := httptest.NewRecorder()

		handler.VerifyOGP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("Method %s: Expected status %d, got %d", method, http.StatusMethodNotAllowed, status)
		}

		if body := rr.Body.String(); body != "Method not allowed\n" {
			t.Errorf("Method %s: Expected 'Method not allowed\\n', got %q", method, body)
		}
	}
}

func TestOGPHandlerXForwardedFor(t *testing.T) {
	handler := handlers.NewOGPHandler()

	// Make 10 requests with X-Forwarded-For header
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/ogp/verify", bytes.NewReader([]byte(`{"url":"https://example.com"}`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "10.0.0.1")
		req.RemoteAddr = "127.0.0.1:12345" // This should be ignored

		rr := httptest.NewRecorder()
		handler.VerifyOGP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Request %d: Expected status %d, got %d", i+1, http.StatusOK, status)
		}
	}

	// 11th request should be rate limited
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ogp/verify", bytes.NewReader([]byte(`{"url":"https://example.com"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	req.RemoteAddr = "127.0.0.1:12345"

	rr := httptest.NewRecorder()
	handler.VerifyOGP(rr, req)

	if status := rr.Code; status != http.StatusTooManyRequests {
		t.Errorf("Expected status %d for rate limited request, got %d", http.StatusTooManyRequests, status)
	}
}

func TestOGPResponseStructure(t *testing.T) {
	handler := handlers.NewOGPHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ogp/verify", bytes.NewReader([]byte(`{"url":"https://example.com"}`)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.VerifyOGP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}

	var resp models.OGPResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check response structure
	if resp.URL == "" {
		t.Error("Expected URL to be set")
	}
	if resp.OGPData.Title == "" && resp.OGPData.Description == "" && resp.OGPData.Image == "" {
		t.Log("OGPData appears to be empty (this may be expected for some sites)")
	}
	
	// Check platform previews exist
	if resp.Previews.Twitter.Platform != "twitter" {
		t.Error("Expected Twitter preview to be set")
	}
	if resp.Previews.Facebook.Platform != "facebook" {
		t.Error("Expected Facebook preview to be set")
	}
	if resp.Previews.Discord.Platform != "discord" {
		t.Error("Expected Discord preview to be set")
	}
	if resp.Timestamp.IsZero() {
		t.Error("Expected Timestamp to be set")
	}

	// Check timestamp is recent
	if time.Since(resp.Timestamp) > 5*time.Second {
		t.Error("Expected timestamp to be recent")
	}
}