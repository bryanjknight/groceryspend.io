variable "region" {
  default = "nyc3"
}

variable "droplet_image" {
  default = "s-1vcpu-1gb"
}

variable "k8s_worker_image" {
  default = "s-2vcpu-2gb"
}

variable "ip_range" {
  default = "10.10.10.0/24"
}