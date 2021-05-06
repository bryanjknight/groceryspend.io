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

resource "null_resource" "setup_db" {
  depends_on = [
    digitalocean_database_db.receipts,
    digitalocean_database_db.users,
    digitalocean_database_user.receipts,
    digitalocean_database_user.users
  ]
  # setup receiptsdb and receipts
  provisioner "local-exec" {

    command = "psql --set=sslmode=require -U ${digitalocean_database_cluster.postgres.user} -p ${digitalocean_database_cluster.postgres.port} -h ${digitalocean_database_cluster.postgres.host} -w receiptsdb -c 'CREATE EXTENSION \"uuid-ossp\"'"

    environment = {
      PGPASSWORD = "${digitalocean_database_cluster.postgres.password}"
    }
  }
  provisioner "local-exec" {
    command = "psql --set=sslmode=require -U ${digitalocean_database_cluster.postgres.user} -p ${digitalocean_database_cluster.postgres.port} -h ${digitalocean_database_cluster.postgres.host} -w receiptsdb -c 'GRANT ALL PRIVILEGES ON DATABASE receiptsdb TO receipts'"

    environment = {
      PGPASSWORD = "${digitalocean_database_cluster.postgres.password}"
    }
  }

  # setup usersdb and users
  provisioner "local-exec" {

    command = "psql --set=sslmode=require -U ${digitalocean_database_cluster.postgres.user} -p ${digitalocean_database_cluster.postgres.port} -h ${digitalocean_database_cluster.postgres.host} -w usersdb -c 'CREATE EXTENSION \"uuid-ossp\"'"

    environment = {
      PGPASSWORD = "${digitalocean_database_cluster.postgres.password}"
    }
  }
  provisioner "local-exec" {
    command = "psql --set=sslmode=require -U ${digitalocean_database_cluster.postgres.user} -p ${digitalocean_database_cluster.postgres.port} -h ${digitalocean_database_cluster.postgres.host} -w usersdb -c 'GRANT ALL PRIVILEGES ON DATABASE usersdb TO users'"

    environment = {
      PGPASSWORD = "${digitalocean_database_cluster.postgres.password}"
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
