resource "digitalocean_vpc" "staging" {
  name     = "groceryspend-staging"
  region   = "${var.region}"
  ip_range = "${var.ip_range}"
}
