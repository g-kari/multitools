package services_test

import (
	"testing"

	"ogp-verification-service/internal/services"
)

func TestOGPServiceIntegration(t *testing.T) {
	service := services.NewOGPService()

	tests := []struct {
		name          string
		url           string
		expectError   bool
		expectTitle   bool
		expectImage   bool
		expectDesc    bool
		minTitleLen   int
		minDescLen    int
	}{
		{
			name:        "GitHub - Complete OGP",
			url:         "https://github.com",
			expectError: false,
			expectTitle: true,
			expectImage: true,
			expectDesc:  true,
			minTitleLen: 10,
			minDescLen:  50,
		},
		{
			name:        "Wikipedia - Good OGP",
			url:         "https://www.wikipedia.org",
			expectError: false,
			expectTitle: true,
			expectImage: true,
			expectDesc:  true,
			minTitleLen: 5,
			minDescLen:  20,
		},
		{
			name:        "Example.com - Minimal OGP",
			url:         "https://example.com",
			expectError: false,
			expectTitle: false,
			expectImage: false,
			expectDesc:  false,
		},
		{
			name:        "Private IP - Should fail",
			url:         "http://192.168.1.1",
			expectError: true,
		},
		{
			name:        "Invalid URL - Should fail",
			url:         "not-a-url",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.FetchOGPData(tt.url)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for URL %s, but got none", tt.url)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error for URL %s: %v", tt.url, err)
			}

			// Check basic response structure
			if resp == nil {
				t.Fatal("Expected response, got nil")
			}

			if resp.URL != tt.url {
				t.Errorf("Expected URL %s, got %s", tt.url, resp.URL)
			}

			// OGPData, Validation, and Previews are structs, not pointers, so they're always present

			// Check OGP data expectations
			if tt.expectTitle && resp.OGPData.Title == "" {
				t.Errorf("Expected title for %s, got empty", tt.url)
			}
			if !tt.expectTitle && resp.OGPData.Title != "" {
				t.Logf("Unexpected title for %s: %s", tt.url, resp.OGPData.Title)
			}

			if tt.expectImage && resp.OGPData.Image == "" {
				t.Errorf("Expected image for %s, got empty", tt.url)
			}
			if !tt.expectImage && resp.OGPData.Image != "" {
				t.Logf("Unexpected image for %s: %s", tt.url, resp.OGPData.Image)
			}

			if tt.expectDesc && resp.OGPData.Description == "" {
				t.Errorf("Expected description for %s, got empty", tt.url)
			}
			if !tt.expectDesc && resp.OGPData.Description != "" {
				t.Logf("Unexpected description for %s: %s", tt.url, resp.OGPData.Description)
			}

			// Check minimum lengths
			if tt.expectTitle && len(resp.OGPData.Title) < tt.minTitleLen {
				t.Errorf("Expected title length >= %d for %s, got %d", tt.minTitleLen, tt.url, len(resp.OGPData.Title))
			}

			if tt.expectDesc && len(resp.OGPData.Description) < tt.minDescLen {
				t.Errorf("Expected description length >= %d for %s, got %d", tt.minDescLen, tt.url, len(resp.OGPData.Description))
			}

			// Check validation results
			if tt.expectTitle != resp.Validation.Checks.HasTitle {
				t.Errorf("Expected HasTitle %v for %s, got %v", tt.expectTitle, tt.url, resp.Validation.Checks.HasTitle)
			}

			if tt.expectImage != resp.Validation.Checks.HasImage {
				t.Errorf("Expected HasImage %v for %s, got %v", tt.expectImage, tt.url, resp.Validation.Checks.HasImage)
			}

			if tt.expectDesc != resp.Validation.Checks.HasDescription {
				t.Errorf("Expected HasDescription %v for %s, got %v", tt.expectDesc, tt.url, resp.Validation.Checks.HasDescription)
			}

			// Check platform previews
			if resp.Previews.Twitter.Platform != "twitter" {
				t.Errorf("Expected Twitter platform to be 'twitter', got %s", resp.Previews.Twitter.Platform)
			}
			if resp.Previews.Twitter.MaxTitleLen != 70 {
				t.Errorf("Expected Twitter max title length to be 70, got %d", resp.Previews.Twitter.MaxTitleLen)
			}
			if resp.Previews.Twitter.MaxDescLen != 200 {
				t.Errorf("Expected Twitter max desc length to be 200, got %d", resp.Previews.Twitter.MaxDescLen)
			}

			if resp.Previews.Facebook.Platform != "facebook" {
				t.Errorf("Expected Facebook platform to be 'facebook', got %s", resp.Previews.Facebook.Platform)
			}
			if resp.Previews.Facebook.MaxTitleLen != 100 {
				t.Errorf("Expected Facebook max title length to be 100, got %d", resp.Previews.Facebook.MaxTitleLen)
			}
			if resp.Previews.Facebook.MaxDescLen != 300 {
				t.Errorf("Expected Facebook max desc length to be 300, got %d", resp.Previews.Facebook.MaxDescLen)
			}

			if resp.Previews.Discord.Platform != "discord" {
				t.Errorf("Expected Discord platform to be 'discord', got %s", resp.Previews.Discord.Platform)
			}
			if resp.Previews.Discord.MaxTitleLen != 256 {
				t.Errorf("Expected Discord max title length to be 256, got %d", resp.Previews.Discord.MaxTitleLen)
			}
			if resp.Previews.Discord.MaxDescLen != 2048 {
				t.Errorf("Expected Discord max desc length to be 2048, got %d", resp.Previews.Discord.MaxDescLen)
			}

			// Check timestamp
			if resp.Timestamp.IsZero() {
				t.Error("Expected timestamp to be set")
			}
		})
	}
}

func TestOGPServiceValidation(t *testing.T) {
	service := services.NewOGPService()

	tests := []struct {
		name           string
		url            string
		expectWarnings bool
		expectErrors   bool
	}{
		{
			name:           "Complete OGP - GitHub",
			url:            "https://github.com",
			expectWarnings: false,
			expectErrors:   false,
		},
		{
			name:           "Incomplete OGP - Example",
			url:            "https://example.com",
			expectWarnings: true,
			expectErrors:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.FetchOGPData(tt.url)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			hasWarnings := len(resp.Validation.Warnings) > 0
			hasErrors := len(resp.Validation.Errors) > 0

			if tt.expectWarnings != hasWarnings {
				t.Errorf("Expected warnings %v for %s, got %v (warnings: %v)", tt.expectWarnings, tt.url, hasWarnings, resp.Validation.Warnings)
			}

			if tt.expectErrors != hasErrors {
				t.Errorf("Expected errors %v for %s, got %v (errors: %v)", tt.expectErrors, tt.url, hasErrors, resp.Validation.Errors)
			}

			// Overall validation should be true unless there are errors
			expectedValid := !hasErrors
			if resp.Validation.IsValid != expectedValid {
				t.Errorf("Expected IsValid %v for %s, got %v", expectedValid, tt.url, resp.Validation.IsValid)
			}
		})
	}
}

func TestOGPServiceCharacterLimits(t *testing.T) {
	service := services.NewOGPService()

	// Test with Amazon which has long descriptions
	resp, err := service.FetchOGPData("https://www.amazon.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(resp.Previews.Twitter.Description) > 200 {
		// Should be truncated for Twitter
		expected := len(resp.Previews.Twitter.Warnings) > 0
		if !expected {
			t.Error("Expected warning for Twitter description length, but got none")
		}
	}

	if len(resp.Previews.Facebook.Description) > 300 {
		// Should be truncated for Facebook
		expected := len(resp.Previews.Facebook.Warnings) > 0
		if !expected {
			t.Error("Expected warning for Facebook description length, but got none")
		}
	}

	// Discord should accept longer descriptions without warnings
	if len(resp.OGPData.Description) < 2048 {
		// Discord should not have warnings for reasonable lengths
		if len(resp.Previews.Discord.Warnings) > 0 {
			t.Errorf("Unexpected warnings for Discord: %v", resp.Previews.Discord.Warnings)
		}
	}
}