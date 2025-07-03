package services

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"ogp-verification-service/internal/models"
	"golang.org/x/net/html"
)

type OGPService struct {
	client *http.Client
}

func NewOGPService() *OGPService {
	return &OGPService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *OGPService) FetchOGPData(targetURL string) (*models.OGPResponse, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if s.isPrivateIP(parsedURL.Hostname()) {
		return nil, fmt.Errorf("private IP addresses are not allowed")
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "OGP-Verification-Service/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	ogpData := s.parseOGPTags(string(body))
	validation := s.validateOGPData(ogpData)
	previews := s.generatePlatformPreviews(ogpData)

	return &models.OGPResponse{
		URL:        targetURL,
		OGPData:    ogpData,
		Validation: validation,
		Previews:   previews,
		Timestamp:  time.Now(),
	}, nil
}

func (s *OGPService) parseOGPTags(htmlContent string) models.OGPData {
	ogpData := models.OGPData{}
	
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return ogpData
	}

	s.extractOGPTags(doc, &ogpData)
	return ogpData
}

func (s *OGPService) extractOGPTags(n *html.Node, ogpData *models.OGPData) {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var property, content string
		for _, attr := range n.Attr {
			if attr.Key == "property" {
				property = attr.Val
			} else if attr.Key == "content" {
				content = attr.Val
			}
		}

		switch property {
		case "og:title":
			ogpData.Title = content
		case "og:description":
			ogpData.Description = content
		case "og:image":
			ogpData.Image = content
		case "og:url":
			ogpData.URL = content
		case "og:type":
			ogpData.Type = content
		case "og:site_name":
			ogpData.SiteName = content
		case "og:image:width":
			ogpData.ImageWidth = content
		case "og:image:height":
			ogpData.ImageHeight = content
		case "og:image:alt":
			ogpData.ImageAlt = content
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.extractOGPTags(c, ogpData)
	}
}

func (s *OGPService) validateOGPData(ogpData models.OGPData) models.ValidationResult {
	result := models.ValidationResult{
		IsValid:  true,
		Warnings: []string{},
		Errors:   []string{},
		Checks: models.ValidationChecks{
			HasTitle:       ogpData.Title != "",
			HasDescription: ogpData.Description != "",
			HasImage:       ogpData.Image != "",
			URLValid:       ogpData.URL != "",
		},
	}

	if !result.Checks.HasTitle {
		result.Warnings = append(result.Warnings, "Missing og:title tag")
	}
	if !result.Checks.HasDescription {
		result.Warnings = append(result.Warnings, "Missing og:description tag")
	}
	if !result.Checks.HasImage {
		result.Warnings = append(result.Warnings, "Missing og:image tag")
	}

	if ogpData.Image != "" {
		result.Checks.ImageValid = s.validateImageURL(ogpData.Image)
		if !result.Checks.ImageValid {
			result.Errors = append(result.Errors, "Invalid image URL")
		}
	}

	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result
}

func (s *OGPService) validateImageURL(imageURL string) bool {
	_, err := url.Parse(imageURL)
	return err == nil
}

func (s *OGPService) generatePlatformPreviews(ogpData models.OGPData) models.PlatformPreviews {
	return models.PlatformPreviews{
		Twitter:  s.generateTwitterPreview(ogpData),
		Facebook: s.generateFacebookPreview(ogpData),
		Discord:  s.generateDiscordPreview(ogpData),
	}
}

func (s *OGPService) generateTwitterPreview(ogpData models.OGPData) models.PlatformPreview {
	maxTitleLen := 70
	maxDescLen := 200
	
	preview := models.PlatformPreview{
		Platform:    "twitter",
		Title:       s.truncateString(ogpData.Title, maxTitleLen),
		Description: s.truncateString(ogpData.Description, maxDescLen),
		Image:       ogpData.Image,
		MaxTitleLen: maxTitleLen,
		MaxDescLen:  maxDescLen,
		TitleLength: len(ogpData.Title),
		DescLength:  len(ogpData.Description),
		IsValid:     true,
		Warnings:    []string{},
	}

	if preview.TitleLength > maxTitleLen {
		preview.Warnings = append(preview.Warnings, "Title exceeds Twitter limit (70 characters)")
	}
	if preview.DescLength > maxDescLen {
		preview.Warnings = append(preview.Warnings, "Description exceeds Twitter limit (200 characters)")
	}

	return preview
}

func (s *OGPService) generateFacebookPreview(ogpData models.OGPData) models.PlatformPreview {
	maxTitleLen := 100
	maxDescLen := 300
	
	preview := models.PlatformPreview{
		Platform:    "facebook",
		Title:       s.truncateString(ogpData.Title, maxTitleLen),
		Description: s.truncateString(ogpData.Description, maxDescLen),
		Image:       ogpData.Image,
		MaxTitleLen: maxTitleLen,
		MaxDescLen:  maxDescLen,
		TitleLength: len(ogpData.Title),
		DescLength:  len(ogpData.Description),
		IsValid:     true,
		Warnings:    []string{},
	}

	if preview.TitleLength > maxTitleLen {
		preview.Warnings = append(preview.Warnings, "Title exceeds Facebook limit (100 characters)")
	}
	if preview.DescLength > maxDescLen {
		preview.Warnings = append(preview.Warnings, "Description exceeds Facebook limit (300 characters)")
	}

	return preview
}

func (s *OGPService) generateDiscordPreview(ogpData models.OGPData) models.PlatformPreview {
	maxTitleLen := 256
	maxDescLen := 2048
	
	preview := models.PlatformPreview{
		Platform:    "discord",
		Title:       s.truncateString(ogpData.Title, maxTitleLen),
		Description: s.truncateString(ogpData.Description, maxDescLen),
		Image:       ogpData.Image,
		MaxTitleLen: maxTitleLen,
		MaxDescLen:  maxDescLen,
		TitleLength: len(ogpData.Title),
		DescLength:  len(ogpData.Description),
		IsValid:     true,
		Warnings:    []string{},
	}

	if preview.TitleLength > maxTitleLen {
		preview.Warnings = append(preview.Warnings, "Title exceeds Discord limit (256 characters)")
	}
	if preview.DescLength > maxDescLen {
		preview.Warnings = append(preview.Warnings, "Description exceeds Discord limit (2048 characters)")
	}

	return preview
}

func (s *OGPService) truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen-3] + "..."
}

func (s *OGPService) isPrivateIP(host string) bool {
	privateIPRegex := regexp.MustCompile(`^(10\.|172\.(1[6-9]|2[0-9]|3[01])\.|192\.168\.|127\.|::1|localhost)`)
	return privateIPRegex.MatchString(host)
}