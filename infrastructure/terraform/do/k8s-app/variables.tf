variable "receipts_password" {
  type = string
}
variable "users_password" {
  type = string
}

variable "pg_host" {
  type = string
}

variable "pg_port" {
  type = string
}

variable "cluster_id" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "cluster_endpoint" {
}

variable "cluster_kubeconfig_raw_config" {
}

variable "cluster_kubeconfig_token" {
}

variable "cluster_kubeconfig_ca_cert" {
}

variable "write_kubeconfig" {
  type    = bool
  default = false
}
