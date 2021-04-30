output "bastion_ip" {
  value = digitalocean_droplet.bastion.*.ipv4_address
}