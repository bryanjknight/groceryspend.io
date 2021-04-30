output "bastion_ip" {
  value = digitalocean_droplet.bastion.*.ipv4_address
}

# output "db_uri" {
#   value = digitalocean_database_cluster.postgres.private_uri
# }

# output "db_admin_user" {
#   value = digitalocean_database_cluster.postgres.user
# }

# output "db_admin_password" {
#   value = digitalocean_database_cluster.postgres.password
# }

# output "receipts_password" {
#   value = digitalocean_database_user.receipts.password
# }

# output "users_password" {
#   value = digitalocean_database_user.users.password
# }