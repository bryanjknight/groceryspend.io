resource "digitalocean_database_cluster" "postgres" {
  name                 = "postgres-cluster"
  engine               = "pg"
  version              = "12"
  size                 = "db-s-1vcpu-1gb"
  region               = var.region
  node_count           = 1
  private_network_uuid = digitalocean_vpc.groceryspend-vpc.id
}

resource "digitalocean_database_db" "receipts" {
  depends_on = [
    digitalocean_database_cluster.postgres
  ]

  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "receiptsdb"
}

resource "digitalocean_database_user" "receipts" {
  depends_on = [
    digitalocean_database_cluster.postgres
  ]

  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "receipts"
}

resource "digitalocean_database_db" "users" {
  depends_on = [
    digitalocean_database_cluster.postgres
  ]

  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "usersdb"
}

resource "digitalocean_database_user" "users" {
  depends_on = [
    digitalocean_database_cluster.postgres
  ]

  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "users"
}

resource "digitalocean_database_db" "payments" {
  depends_on = [
    digitalocean_database_cluster.postgres
  ]

  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "paymentsdb"
}

resource "digitalocean_database_user" "payments" {
  depends_on = [
    digitalocean_database_cluster.postgres
  ]

  cluster_id = digitalocean_database_cluster.postgres.id
  name       = "payments"
}

resource "null_resource" "setup_db" {
  depends_on = [
    digitalocean_database_db.receipts,
    digitalocean_database_db.users,
    digitalocean_database_db.payments,
    digitalocean_database_user.receipts,
    digitalocean_database_user.users,
    digitalocean_database_user.payments
  ]

  # setup receiptsdb and receipts via the bastion
  provisioner "local-exec" {

    command = "../../server/scripts/init-db.sh"

    environment = {
      POSTGRES_USER = "${digitalocean_database_cluster.postgres.user}"
      PGPASSWORD = "${digitalocean_database_cluster.postgres.password}"
      PGHOST =  "${digitalocean_database_cluster.postgres.host}"
      PGPORT = "${digitalocean_database_cluster.postgres.port}"
      POSTGRES_DB = "defaultdb" # digitalocean's default database
    }
  }
}

resource "digitalocean_database_firewall" "only_vpc_traffic_fw" {
  cluster_id = digitalocean_database_cluster.postgres.id

  depends_on = [
    null_resource.setup_db
  ]

  rule {
    type  = "tag"
    value = "k8s"
  }

  rule {
    type  = "tag"
    value = "bastion"
  }
}
