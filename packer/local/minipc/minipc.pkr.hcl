packer {
  required_plugins {
    qemu = {
      source  = "github.com/hashicorp/qemu"
      version = ">= 1.1.0"
    }
    shell = {
      source  = "github.com/hashicorp/shell"
      version = ">= 1.0.0"
    }
  }
}

locals {
  timestamp      = formatdate("YYYYMMDD-hhmmss", timestamp())
  image_name     = "${var.snapshot_name}-${local.timestamp}"
}

source "qemu" "minipc" {
  iso_url          = var.iso_url
  iso_checksum     = var.iso_checksum
  output_directory = "output-qemu"
  vm_name          = "minipc.qcow2"
  disk_size        = var.disk_size
  memory           = var.memory
  cores            = var.cores
  format           = "qcow2"
  ssh_username     = "root"
  ssh_password     = "packer"
  ssh_timeout      = "30m"
  boot_wait        = "10s"
  boot_command = [
    "<esc><esc><esc><esc>e<wait>",
    "linux /casper/vmlinuz ",
    "autoinstall ds=nocloud-net;s=http://{{ .HTTPIP }}:{{ .HTTPPort }}/user-data --- ",
    "<f10>"
  ]
  http_directory   = "http"
  shutdown_command = "sudo shutdown -P now"
}

build {
  sources = ["source.qemu.minipc"]

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
      image_name     = local.image_name
      franky_version = var.franky_version
    }
  }
}