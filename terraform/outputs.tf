output "server_ip" {
  description = "IP address of the OGP verification server"
  value       = sakuracloud_server.ogp_server.ip_address
}

output "server_id" {
  description = "ID of the OGP verification server"
  value       = sakuracloud_server.ogp_server.id
}

output "api_endpoint" {
  description = "API endpoint URL"
  value       = "https://api.${var.domain_name}"
}

output "frontend_url" {
  description = "Frontend URL"
  value       = "https://${var.domain_name}"
}

output "ssh_connection" {
  description = "SSH connection command"
  value       = "ssh ubuntu@${sakuracloud_server.ogp_server.ip_address}"
}

output "cloudflare_zone_id" {
  description = "Cloudflare Zone ID"
  value       = var.cloudflare_zone_id
}

output "dns_records" {
  description = "DNS records created"
  value = {
    api  = cloudflare_record.ogp_api.hostname
    www  = cloudflare_record.ogp_www.hostname
    root = cloudflare_record.ogp_root.hostname
  }
}