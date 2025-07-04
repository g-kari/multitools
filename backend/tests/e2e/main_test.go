package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	defaultAPIBaseURL = "http://localhost:8080"
)

type OGPRequest struct {
	URL string `json:"url"`
}

type OGPResponse struct {
	URL        string            `json:"url"`
	OGPData    OGPData          `json:"ogp_data"`
	Validation ValidationResult `json:"validation"`
	Previews   PlatformPreviews `json:"previews"`
	Timestamp  time.Time        `json:"timestamp"`
}

type OGPData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	SiteName    string `json:"site_name"`
	ImageWidth  string `json:"image_width"`
	ImageHeight string `json:"image_height"`
	ImageAlt    string `json:"image_alt"`
}

type ValidationResult struct {
	IsValid  bool             `json:"is_valid"`
	Warnings []string         `json:"warnings"`
	Errors   []string         `json:"errors"`
	Checks   ValidationChecks `json:"checks"`
}

type ValidationChecks struct {
	HasTitle       bool `json:"has_title"`
	HasDescription bool `json:"has_description"`
	HasImage       bool `json:"has_image"`
	ImageValid     bool `json:"image_valid"`
	URLValid       bool `json:"url_valid"`
}

type PlatformPreviews struct {
	Twitter  PlatformPreview `json:"twitter"`
	Facebook PlatformPreview `json:"facebook"`
	Discord  PlatformPreview `json:"discord"`
}

type PlatformPreview struct {
	Platform     string   `json:"platform"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Image        string   `json:"image"`
	IsValid      bool     `json:"is_valid"`
	Warnings     []string `json:"warnings"`
	TitleLength  int      `json:"title_length"`
	DescLength   int      `json:"desc_length"`
	MaxTitleLen  int      `json:"max_title_len"`
	MaxDescLen   int      `json:"max_desc_len"`
}

func getAPIBaseURL() string {
	if url := os.Getenv("API_BASE_URL"); url != "" {
		return url
	}
	return defaultAPIBaseURL
}

func makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	baseURL := getAPIBaseURL()
	url := baseURL + endpoint

	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	return client.Do(req)
}

func TestHealthEndpoint(t *testing.T) {
	resp, err := makeRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to make health check request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestOGPVerifyEndpoint_GitHub(t *testing.T) {
	req := OGPRequest{URL: "https://github.com"}
	
	resp, err := makeRequest("POST", "/api/v1/ogp/verify", req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var ogpResp OGPResponse
	if err := json.NewDecoder(resp.Body).Decode(&ogpResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response structure
	if ogpResp.URL != "https://github.com" {
		t.Errorf("Expected URL to be https://github.com, got %s", ogpResp.URL)
	}

	if ogpResp.OGPData.Title == "" {
		t.Error("Expected GitHub to have a title")
	}

	if ogpResp.OGPData.Description == "" {
		t.Error("Expected GitHub to have a description")
	}

	if ogpResp.OGPData.Image == "" {
		t.Error("Expected GitHub to have an image")
	}

	// Verify platform previews
	if ogpResp.Previews.Twitter.Platform != "twitter" {
		t.Errorf("Expected Twitter platform, got %s", ogpResp.Previews.Twitter.Platform)
	}

	if ogpResp.Previews.Facebook.Platform != "facebook" {
		t.Errorf("Expected Facebook platform, got %s", ogpResp.Previews.Facebook.Platform)
	}

	if ogpResp.Previews.Discord.Platform != "discord" {
		t.Errorf("Expected Discord platform, got %s", ogpResp.Previews.Discord.Platform)
	}

	// Verify validation results
	if !ogpResp.Validation.IsValid {
		t.Error("Expected GitHub OGP to be valid")
	}

	if len(ogpResp.Validation.Errors) > 0 {
		t.Errorf("Expected no errors, got %v", ogpResp.Validation.Errors)
	}
}

func TestOGPVerifyEndpoint_Wikipedia(t *testing.T) {
	req := OGPRequest{URL: "https://www.wikipedia.org"}
	
	resp, err := makeRequest("POST", "/api/v1/ogp/verify", req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var ogpResp OGPResponse
	if err := json.NewDecoder(resp.Body).Decode(&ogpResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if ogpResp.OGPData.Title == "" {
		t.Error("Expected Wikipedia to have a title")
	}

	if ogpResp.Validation.Checks.HasTitle != true {
		t.Error("Expected Wikipedia to have title check as true")
	}
}

func TestOGPVerifyEndpoint_ExampleCom(t *testing.T) {
	req := OGPRequest{URL: "https://example.com"}
	
	resp, err := makeRequest("POST", "/api/v1/ogp/verify", req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var ogpResp OGPResponse
	if err := json.NewDecoder(resp.Body).Decode(&ogpResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Example.com typically has minimal OGP data
	if len(ogpResp.Validation.Warnings) == 0 {
		t.Log("Example.com unexpectedly has complete OGP data")
	}

	// Should still be valid overall
	if !ogpResp.Validation.IsValid {
		t.Error("Expected response to be valid even with warnings")
	}
}

func TestOGPVerifyEndpoint_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		request        interface{}
		expectedStatus int
	}{
		{
			name:           "Empty URL",
			request:        OGPRequest{URL: ""},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON",
			request:        `{"invalid": json}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Private IP",
			request:        OGPRequest{URL: "http://192.168.1.1"},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			if str, ok := tt.request.(string); ok {
				// Invalid JSON case
				req, _ := http.NewRequest("POST", getAPIBaseURL()+"/api/v1/ogp/verify", bytes.NewReader([]byte(str)))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{Timeout: 30 * time.Second}
				resp, err = client.Do(req)
			} else {
				resp, err = makeRequest("POST", "/api/v1/ogp/verify", tt.request)
			}

			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestRateLimiting(t *testing.T) {
	req := OGPRequest{URL: "https://example.com"}

	// Make 10 requests (should all succeed)
	for i := 0; i < 10; i++ {
		resp, err := makeRequest("POST", "/api/v1/ogp/verify", req)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i+1, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i+1, resp.StatusCode)
		}
	}

	// 11th request should be rate limited
	resp, err := makeRequest("POST", "/api/v1/ogp/verify", req)
	if err != nil {
		t.Fatalf("Rate limit test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 for rate limited request, got %d", resp.StatusCode)
	}
}

func TestCORSHeaders(t *testing.T) {
	resp, err := makeRequest("OPTIONS", "/api/v1/ogp/verify", nil)
	if err != nil {
		t.Fatalf("Failed to make OPTIONS request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS request, got %d", resp.StatusCode)
	}

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}

	for header, expected := range expectedHeaders {
		if value := resp.Header.Get(header); value != expected {
			t.Errorf("Expected header %s to be '%s', got '%s'", header, expected, value)
		}
	}
}

func TestMethodNotAllowed(t *testing.T) {
	methods := []string{"GET", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		resp, err := makeRequest(method, "/api/v1/ogp/verify", nil)
		if err != nil {
			t.Fatalf("Failed to make %s request: %v", method, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Method %s: Expected status 405, got %d", method, resp.StatusCode)
		}
	}
}

func TestAPILatency(t *testing.T) {
	req := OGPRequest{URL: "https://github.com"}
	
	start := time.Now()
	resp, err := makeRequest("POST", "/api/v1/ogp/verify", req)
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// API should respond within 10 seconds (as per requirements)
	if duration > 10*time.Second {
		t.Errorf("API response took %v, expected < 10s", duration)
	}

	t.Logf("API response time: %v", duration)
}

func TestConcurrentRequests(t *testing.T) {
	const numRequests = 5
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			req := OGPRequest{URL: fmt.Sprintf("https://example.com?id=%d", id)}
			resp, err := makeRequest("POST", "/api/v1/ogp/verify", req)
			if err != nil {
				results <- fmt.Errorf("request %d failed: %v", id, err)
				return
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("request %d: expected status 200, got %d", id, resp.StatusCode)
				return
			}

			results <- nil
		}(i)
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		if err := <-results; err != nil {
			t.Error(err)
		}
	}
}