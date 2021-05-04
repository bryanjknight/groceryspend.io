output "vpc_id" {
  value = digitalocean_vpc.groceryspend-vpc.id
}

output "bastion_ipv4_address" {
  value = digitalocean_droplet.bastion.*.ipv4_address
}