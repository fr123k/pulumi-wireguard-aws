// Path to the OVA to boot (base Ubuntu or pre-baked WireGuard OVA)
variable "source_path" {
  type        = string
  description = "Path to the OVA file to boot as the test VM"
}

// Target OS/arch for the verify binary build (GOOS/GOARCH).
variable "target_goos" {
  type        = string
  default     = "linux"
  description = "Target OS for the verify binary (GOOS)"
}

variable "target_goarch" {
  type        = string
  default     = "amd64"
  description = "Target architecture for the verify binary (GOARCH)"
}

variable "home_dir" {
  type    = string
  default = env("HOME")
}

// SSH credentials (must match the image's pre-configured user)
variable "ssh_username" {
  type    = string
  default = "packer"
}

variable "ssh_password" {
  type    = string
  default = "packer"
}

variable "ssh_private_key_file" {
  type    = string
  default = "${env("HOME")}/.virtualbox/packer/packer_private_key_file"
}

// Cloud-init template to test
variable "cloud_init_file" {
  type        = string
  description = "Path to the cloud-init template file (e.g. cloud-init/wireguard-prebaked.txt)"
}

// Verify target and mode
variable "target" {
  type        = string
  default     = "wireguard"
  description = "Target to verify: wireguard, temporal, franky, etc."
}

variable "mode" {
  type        = string
  default     = "deployed"
  description = "Verify mode: prebaked or deployed"
}

// Template variables rendered into the cloud-init script
variable "domain" {
  type    = string
  default = "wg.test.local"
}

variable "temporal_domain" {
  type    = string
  default = "temporal.dunebot.io"
}

variable "dunebot_domain" {
  type    = string
  default = "githubapp.dunebot.io"
}

variable "client_publickey" {
  type    = string
  default = "test-client-public-key-placeholder"
}

variable "client_ip_address" {
  type    = string
  default = "172.16.16.2"
}

variable "secret_operator_token" {
  type    = string
  default = ""
  sensitive = true
}

// VM resources
variable "vm_name" {
  type    = string
  default = "cloudinit-test"
}

variable "memory" {
  type    = number
  default = 2048
}

variable "cpus" {
  type    = number
  default = 2
}

variable "output_directory" {
  type    = string
  default = "packer/cloudinit-test/output"
}
