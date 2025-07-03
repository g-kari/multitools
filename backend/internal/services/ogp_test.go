package services

import (
	"testing"
	"ogp-verification-service/internal/models"
)

func TestOGPService_parseOGPTags(t *testing.T) {
	service := NewOGPService()
	
	tests := []struct {
		name     string
		html     string
		expected models.OGPData
	}{
		{
			name: "Complete OGP tags",
			html: `
				<html>
					<head>
						<meta property="og:title" content="Test Title" />
						<meta property="og:description" content="Test Description" />
						<meta property="og:image" content="https://example.com/image.jpg" />
						<meta property="og:url" content="https://example.com" />
						<meta property="og:type" content="website" />
						<meta property="og:site_name" content="Test Site" />
					</head>
				</html>
			`,
			expected: models.OGPData{
				Title:       "Test Title",
				Description: "Test Description",
				Image:       "https://example.com/image.jpg",
				URL:         "https://example.com",
				Type:        "website",
				SiteName:    "Test Site",
			},
		},
		{
			name: "Partial OGP tags",
			html: `
				<html>
					<head>
						<meta property="og:title" content="Partial Title" />
						<meta property="og:description" content="Partial Description" />
					</head>
				</html>
			`,
			expected: models.OGPData{
				Title:       "Partial Title",
				Description: "Partial Description",
			},
		},
		{
			name: "No OGP tags",
			html: `
				<html>
					<head>
						<title>Regular Title</title>
					</head>
				</html>
			`,
			expected: models.OGPData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.parseOGPTags(tt.html)
			
			if result.Title != tt.expected.Title {
				t.Errorf("Expected title %s, got %s", tt.expected.Title, result.Title)
			}
			if result.Description != tt.expected.Description {
				t.Errorf("Expected description %s, got %s", tt.expected.Description, result.Description)
			}
			if result.Image != tt.expected.Image {
				t.Errorf("Expected image %s, got %s", tt.expected.Image, result.Image)
			}
			if result.URL != tt.expected.URL {
				t.Errorf("Expected URL %s, got %s", tt.expected.URL, result.URL)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Expected type %s, got %s", tt.expected.Type, result.Type)
			}
			if result.SiteName != tt.expected.SiteName {
				t.Errorf("Expected site name %s, got %s", tt.expected.SiteName, result.SiteName)
			}
		})
	}
}

func TestOGPService_validateOGPData(t *testing.T) {
	service := NewOGPService()
	
	tests := []struct {
		name     string
		data     models.OGPData
		expected bool
	}{
		{
			name: "Valid OGP data",
			data: models.OGPData{
				Title:       "Test Title",
				Description: "Test Description",
				Image:       "https://example.com/image.jpg",
				URL:         "https://example.com",
			},
			expected: true,
		},
		{
			name: "Invalid image URL",
			data: models.OGPData{
				Title:       "Test Title",
				Description: "Test Description",
				Image:       "invalid-url",
				URL:         "https://example.com",
			},
			expected: false,
		},
		{
			name: "Missing required fields",
			data: models.OGPData{
				URL: "https://example.com",
			},
			expected: true, // Only warnings, no errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.validateOGPData(tt.data)
			
			if result.IsValid != tt.expected {
				t.Errorf("Expected IsValid %v, got %v", tt.expected, result.IsValid)
			}
		})
	}
}

func TestOGPService_generateTwitterPreview(t *testing.T) {
	service := NewOGPService()
	
	data := models.OGPData{
		Title:       "This is a very long title that exceeds Twitter's character limit",
		Description: "This is a description",
		Image:       "https://example.com/image.jpg",
	}
	
	result := service.generateTwitterPreview(data)
	
	if result.Platform != "twitter" {
		t.Errorf("Expected platform twitter, got %s", result.Platform)
	}
	
	if result.MaxTitleLen != 70 {
		t.Errorf("Expected max title length 70, got %d", result.MaxTitleLen)
	}
	
	if result.MaxDescLen != 200 {
		t.Errorf("Expected max description length 200, got %d", result.MaxDescLen)
	}
	
	if len(result.Warnings) == 0 {
		t.Error("Expected warnings for long title")
	}
}

func TestOGPService_isPrivateIP(t *testing.T) {
	service := NewOGPService()
	
	tests := []struct {
		host     string
		expected bool
	}{
		{"127.0.0.1", true},
		{"localhost", true},
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"8.8.8.8", false},
		{"google.com", false},
		{"example.com", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			result := service.isPrivateIP(tt.host)
			if result != tt.expected {
				t.Errorf("Expected %v for %s, got %v", tt.expected, tt.host, result)
			}
		})
	}
}

func TestOGPService_truncateString(t *testing.T) {
	service := NewOGPService()
	
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is..."},
		{"exactly10chars", 10, "exactly10chars"},
		{"", 10, ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := service.truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}