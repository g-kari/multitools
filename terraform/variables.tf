variable "sakuracloud_token" {
  description = "Sakura Cloud API Token"
  type        = string
  sensitive   = true
}

variable "sakuracloud_secret" {
  description = "Sakura Cloud API Secret"
  type        = string
  sensitive   = true
}

variable "sakuracloud_zone" {
  description = "Sakura Cloud Zone"
  type        = string
  default     = "tk1a"
}

variable "cloudflare_api_token" {
  description = "Cloudflare API Token"
  type        = string
  sensitive   = true
}

variable "cloudflare_zone_id" {
  description = "Cloudflare Zone ID"
  type        = string
}

variable "domain_name" {
  description = "Domain name for the service"
  type        = string
}

variable "create_cloudflare_zone" {
  description = "Whether to create a new Cloudflare zone"
  type        = bool
  default     = false
}

variable "server_password" {
  description = "Password for the server"
  type        = string
  sensitive   = true
}

variable "ssh_public_key" {
  description = "SSH Public Key for server access"
  type        = string
}

variable "environment" {
  description = "Environment name (production, staging, etc.)"
  type        = string
  default     = "production"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "ogp-verification"
}