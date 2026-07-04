packer {
  required_plugins {
    hcloud = {
      source  = "github.com/hetznercloud/hcloud"
      version = ">= 1.3.0"
    }
  }
}

locals {
  timestamp = formatdate("YYYYMMDD-hhmmss", timestamp())
  snapshot_name = "${var.snapshot_name}-${local.timestamp}"
}

source "hcloud" "franky" {
  image        = var.base_image
  location     = var.location
  server_type  = var.server_type
  server_name  = "packer-franky-builder"
  ssh_username = "root"

  snapshot_name = local.snapshot_name
  snapshot_labels = {
    "app"                      = "franky"
    "build_timestamp"          = local.timestamp
    "franky_version"           = var.franky_version
    "packer_build"             = "true"
  }
}

build {
  sources = ["source.hcloud.franky"]

  # Generate versions.env file from template
  provisioner "file" {
    content = templatefile("${path.root}/scripts/versions.env.tpl", {
      franky_version = var.franky_version
    })
    destination = "/tmp/versions.env"
  }

  # Upload all provisioner scripts
  provisioner "file" {
    source      = "${path.root}/scripts/"
    destination = "/tmp/"
  }

  # Make scripts executable and run them in order
  provisioner "shell" {
    inline = [
      "chmod +x /tmp/*.sh",
      ". /tmp/versions.env && /tmp/01-base-packages.sh",
      ". /tmp/versions.env && /tmp/02-franky-binaries.sh",
      "/tmp/03-nginx-setup.sh",
      "/tmp/04-systemd-services.sh",
      "/tmp/05-security-hardening.sh",
      "/tmp/06-cleanup.sh"
    ]
    environment_vars = [
      "DEBIAN_FRONTEND=noninteractive"
    ]
  }

  post-processor "manifest" {
    output     = "manifest.json"
    strip_path = true
    custom_data = {
      snapshot_name      = local.snapshot_name
      franky_version = var.franky_version
    }
  }
}
