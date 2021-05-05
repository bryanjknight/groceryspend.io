output "vpc_id" {
  value = digitalocean_vpc.groceryspend-vpc.id
}

output "bastion_ipv4_address" {
  value = digitalocean_droplet.bastion.*.ipv4_address
}

output "postgres_host" {
  value = digitalocean_database_cluster.postgres.host
}

output "postgres_port" {
  value = digitalocean_database_cluster.postgres.port
}


output "receipts_password" {
  value = digitalocean_database_user.receipts.password
  sensitive = true
}

output "users_password" {
  value = digitalocean_database_user.users.password
  sensitive = true
}