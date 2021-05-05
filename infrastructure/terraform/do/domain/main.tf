terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.4.0"
    }
  }
}

variable "load_balancer_ip" {
  type = string
}

data "digitalocean_domain" "groceryspend" {
  name = "groceryspend.io"
}

resource "digitalocean_record" "api" {
  domain = data.digitalocean_domain.groceryspend.name
  type   = "A"
  name   = "api"
  value  = "${var.load_balancer_ip}"
}

resource "digitalocean_record" "www" {
  domain = data.digitalocean_domain.groceryspend.name
  type   = "A"
  name   = "www"
  value  = "${var.load_balancer_ip}"
}