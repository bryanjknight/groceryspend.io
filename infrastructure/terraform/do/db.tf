resource "digitalocean_database_cluster" "postgres" {
  name       = "postgres-cluster"
  engine     = "pg"
  version    = "12"
  size       = "db-s-1vcpu-1gb"
  region     = var.region
  node_count = 1
}

resource "digitalocean_database_firewall" "only_vpc_traffic_fw" {
  cluster_id = digitalocean_database_cluster.postgres.id

  rule {
    type  = "tag"
    value = "k8s"
  }

  rule {
    type  = "tag"
    value = "bastion"
  }
}

resource "digitalocean_database_db" "receipts" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "receiptsdb"
}

resource "digitalocean_database_user" "receipts" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "receipts"
}

resource "digitalocean_database_db" "users" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "usersdb"
}

resource "digitalocean_database_user" "users" {
  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "users"
}