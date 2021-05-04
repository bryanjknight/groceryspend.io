resource "digitalocean_vpc" "groceryspend-vpc" {
  name     = "groceryspend-${var.namespace}"
  region   = "${var.region}"
  ip_range = "${var.ip_range}"
}
