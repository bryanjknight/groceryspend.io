variable "namespace" {
  default = "default"
}
variable "region" {
  default = "nyc3"
}

variable "vpc_id" {
  type = string
}

variable "k8s_version" {
  default = "1.20.2-do.0"
}

variable "k8s_worker_image" {
  default = "s-2vcpu-2gb"
}
variable "k8s_node_count" {
  type    = number
  default = 2
}
