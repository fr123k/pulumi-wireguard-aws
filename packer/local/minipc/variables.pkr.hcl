// Local QEMU / Physical Server Configuration

variable "iso_url" {
  type        = string
  default     = "https://releases.ubuntu.com/24.04/ubuntu-24.04.2-live-server-amd64.iso"
  description = "URL of the Ubuntu ISO for building the image"
}

variable "iso_checksum" {
  type        = string
  default     = "sha256:0f004e629e4c68a6a55f6a7a7047e4f8c750772c4f70d22e2e8e695c9c58a7a0"
  description = "Checksum of the Ubuntu ISO"
}

variable "disk_size" {
  type        = string
  default     = "20G"
  description = "Disk size for the QEMU VM"
}

variable "memory" {
  type        = number
  default     = 4096
  description = "Memory in MB for the QEMU VM"
}

variable "cores" {
  type        = number
  default     = 4
  description = "Number of CPU cores for the QEMU VM"
}

variable "snapshot_name" {
  type        = string
  default     = "minipc-prebaked"
  description = "Name prefix for the image"
}

// Component Version Pins
variable "franky_version" {
  type        = string
  default     = "0.32.1"
  description = "franky version"
}