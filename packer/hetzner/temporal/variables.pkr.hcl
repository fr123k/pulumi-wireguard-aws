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
  default     = "temporal-prebaked"
  description = "Name prefix for the snapshot"
}

// Component Version Pins
variable "temporal_cli_version" {
  type        = string
  default     = "1.6.0"
  description = "Temporal CLI version"
}

variable "temporal_server_version" {
  type        = string
  default     = "1.30.0"
  description = "Temporal Server version"
}

variable "temporal_ui_version" {
  type        = string
  default     = "2.45.0"
  description = "Temporal UI Server version"
}

variable "temporal_worker_version" {
  type        = string
  default     = "0.9.6"
  description = "Temporal Worker version"
}

variable "secret_operator_version" {
  type        = string
  default     = "0.4.6"
  description = "Secret Operator version"
}

variable "oauth2_storage_version" {
  type        = string
  default     = "0.2.1"
  description = "OAuth2 Storage version"
}

variable "dunebot_version" {
  type        = string
  default     = "0.3.6"
  description = "DuneBot version"
}
