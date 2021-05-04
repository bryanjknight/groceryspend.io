variable "namespace" {
  default = "default"
}
variable "region" {
  default = "nyc3"
}

variable "ip_range" {
  default = "10.10.10.0/24"
}

variable "droplet_image" {
  default = "s-1vcpu-1gb"
}

variable "project_id" {
  type = string
}

variable "terraform_public_key" {
  type = string
}

variable "pvt_key" {
  type = string
}
