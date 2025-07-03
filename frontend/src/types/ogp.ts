export interface OGPRequest {
  url: string;
}

export interface OGPResponse {
  url: string;
  ogp_data: OGPData;
  validation: ValidationResult;
  previews: PlatformPreviews;
  timestamp: string;
}

export interface OGPData {
  title: string;
  description: string;
  image: string;
  url: string;
  type: string;
  site_name: string;
  image_width: string;
  image_height: string;
  image_alt: string;
}

export interface ValidationResult {
  is_valid: boolean;
  warnings: string[];
  errors: string[];
  checks: ValidationChecks;
}

export interface ValidationChecks {
  has_title: boolean;
  has_description: boolean;
  has_image: boolean;
  image_valid: boolean;
  url_valid: boolean;
}

export interface PlatformPreviews {
  twitter: PlatformPreview;
  facebook: PlatformPreview;
  discord: PlatformPreview;
}

export interface PlatformPreview {
  platform: string;
  title: string;
  description: string;
  image: string;
  is_valid: boolean;
  warnings: string[];
  title_length: number;
  desc_length: number;
  max_title_len: number;
  max_desc_len: number;
}