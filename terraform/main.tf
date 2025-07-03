terraform {
  required_version = ">= 1.0"
  required_providers {
    sakuracloud = {
      source  = "sacloud/sakuracloud"
      version = "~> 2.25"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

provider "sakuracloud" {
  token  = var.sakuracloud_token
  secret = var.sakuracloud_secret
  zone   = var.sakuracloud_zone
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

resource "sakuracloud_server" "ogp_server" {
  name            = "ogp-verification-server"
  description     = "OGP Verification Service Backend"
  disks           = [sakuracloud_disk.ogp_disk.id]
  core            = 1
  memory          = 1
  interface_driver = "virtio"
  
  network_interface {
    upstream = "shared"
  }

  disk_edit_parameter {
    hostname        = "ogp-server"
    password        = var.server_password
    ssh_key_ids     = [sakuracloud_ssh_key.ogp_key.id]
    
    startup_script = file("${path.module}/startup.sh")
  }

  tags = ["ogp", "production"]
}

resource "sakuracloud_disk" "ogp_disk" {
  name              = "ogp-server-disk"
  description       = "OGP Server Root Disk"
  size              = 20
  source_archive_id = data.sakuracloud_archive.ubuntu.id
  
  tags = ["ogp", "production"]
}

resource "sakuracloud_ssh_key" "ogp_key" {
  name        = "ogp-server-key"
  description = "SSH Key for OGP Server"
  public_key  = var.ssh_public_key
}

data "sakuracloud_archive" "ubuntu" {
  name_selectors = ["Ubuntu Server 22.04"]
}

resource "cloudflare_zone" "ogp_zone" {
  count = var.create_cloudflare_zone ? 1 : 0
  zone  = var.domain_name
}

resource "cloudflare_record" "ogp_api" {
  zone_id = var.cloudflare_zone_id
  name    = "api"
  value   = sakuracloud_server.ogp_server.ip_address
  type    = "A"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "ogp_www" {
  zone_id = var.cloudflare_zone_id
  name    = "www"
  value   = var.domain_name
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "ogp_root" {
  zone_id = var.cloudflare_zone_id
  name    = "@"
  value   = sakuracloud_server.ogp_server.ip_address
  type    = "A"
  ttl     = 1
  proxied = true
}

resource "cloudflare_page_rule" "api_cors" {
  zone_id = var.cloudflare_zone_id
  target  = "api.${var.domain_name}/*"
  priority = 1

  actions {
    always_online = "on"
    browser_cache_ttl = 0
    cache_level = "bypass"
    
    cache_key_fields {
      cookie {
        check_presence = []
        include = []
      }
      header {
        check_presence = []
        exclude = ["origin"]
        include = []
      }
      host {
        resolved = true
      }
      query_string {
        exclude = []
        include = []
      }
      user {
        device_type = false
        geo = false
      }
    }
  }
}

resource "cloudflare_zone_settings_override" "ogp_settings" {
  zone_id = var.cloudflare_zone_id
  
  settings {
    always_online = "on"
    always_use_https = "on"
    automatic_https_rewrites = "on"
    browser_cache_ttl = 14400
    cache_level = "aggressive"
    development_mode = "off"
    min_tls_version = "1.2"
    security_level = "medium"
    ssl = "flexible"
  }
}