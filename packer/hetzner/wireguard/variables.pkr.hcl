// Hetzner Cloud Configuration
variable "base_image" {
  type        = string
  default     = "ubuntu-26.04"
  description = "Base image for the snapshot"
}

// Target OS/Architecture for cross-compilation
variable "target_goos" {
  type        = string
  default     = "linux"
  description = "Target OS for verify binary (GOOS)"
}

variable "target_goarch" {
  type        = string
  default     = "amd64"
  description = "Target architecture for verify binary (GOARCH). Use 'amd64' for cx*/cpx*/ccx* server types, 'arm64' for cax* server types"
}

variable "server_type" {
  type        = string
  default     = "cpx12"
  description = "Hetzner server type for building the image"
}

variable "location" {
  type        = string
  default     = "nbg1"
  description = "Hetzner datacenter location"
}

variable "snapshot_name" {
  type        = string
  default     = "wireguard-prebaked"
  description = "Name prefix for the snapshot"
}

// Component Version Pins
variable "wireguard_ui_version" {
  type        = string
  default     = "0.6.2"
  description = "WireGuard UI version (from github.com/ngoduykhanh/wireguard-ui)"
}

variable "secret_operator_version" {
  type        = string
  default     = "0.6.2"
  description = "Secret Operator version (client and server)"
}

// Secret / domain configuration for baking SSL certs into the image
variable "domain" {
  type        = string
  default     = "wg.fr123k.uk"
  description = "WireGuard domain for SSL certificate installation"
}

variable "secret_operator_token" {
  type        = string
  default     = ""
  sensitive   = true
  description = "Secret Operator authentication token for fetching SSL certs at build time. Leave empty to fetch at runtime instead."
}

variable "secret_operator_host" {
  type        = string
  default     = "https://wg.fr123k.uk:8443/secrets"
  description = "Secret Operator server URL. For VirtualBox NAT builds pointing at a local service on the host, use http://10.0.2.2:8888/secrets"
}

// ────────────────────────────────────────────────
// VirtualBox source configuration (local testing)
// Used only when building with -only=virtualbox-ovf.wireguard
// ────────────────────────────────────────────────
variable "source_path" {
  type        = string
  default     = "${env("HOME")}/.virtualbox/packer/packer_base_ubuntu_26.ova"
  description = "Path to the VirtualBox OVA used as the build source."
}

variable "ssh_username" {
  type        = string
  default     = "packer"
  description = "SSH username for the VirtualBox base image."
}

variable "ssh_password" {
  type        = string
  default     = "packer"
  sensitive   = true
  description = "SSH password for the VirtualBox base image (used for sudo alongside the private key)."
}

variable "ssh_private_key_file" {
  type        = string
  default     = "${env("HOME")}/.virtualbox/packer/packer_private_key_file"
  description = "Path to the SSH private key provisioned into the VirtualBox base image."
}

variable "vm_name" {
  type        = string
  default     = "packer-wireguard"
  description = "Name of the VirtualBox VM created during the build."
}

variable "output_directory" {
  type        = string
  default     = "output-wireguard"
  description = "Directory where the built OVA/.box artifact is written."
}

variable "memory" {
  type        = number
  default     = 2048
  description = "Memory (MB) allocated to the VirtualBox build VM."
}

variable "cpus" {
  type        = number
  default     = 2
  description = "Number of CPUs allocated to the VirtualBox build VM."
}