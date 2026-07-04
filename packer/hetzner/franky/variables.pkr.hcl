// Hetzner Cloud Configuration
variable "base_image" {
  type        = string
  default     = "ubuntu-24.04"
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
  default     = "cx23"
  description = "Hetzner server type for building the image"
}

variable "location" {
  type        = string
  default     = "nbg1"
  description = "Hetzner datacenter location"
}

variable "snapshot_name" {
  type        = string
  default     = "franky-prebaked"
  description = "Name prefix for the snapshot"
}

// Component Version Pins
variable "franky_version" {
  type        = string
  default     = "0.32.1"
  description = "franky version"
}
