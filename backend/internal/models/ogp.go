package models

import "time"

type OGPRequest struct {
	URL string `json:"url" validate:"required,url"`
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
	IsValid  bool                `json:"is_valid"`
	Warnings []string           `json:"warnings"`
	Errors   []string           `json:"errors"`
	Checks   ValidationChecks   `json:"checks"`
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
	Platform     string `json:"platform"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	IsValid      bool   `json:"is_valid"`
	Warnings     []string `json:"warnings"`
	TitleLength  int    `json:"title_length"`
	DescLength   int    `json:"desc_length"`
	MaxTitleLen  int    `json:"max_title_len"`
	MaxDescLen   int    `json:"max_desc_len"`
}