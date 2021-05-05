resource "kubernetes_secret" "server-auth" {
  metadata {
    name = "server-auth"
    namespace = kubernetes_namespace.app.metadata.0.name
  }

  data = {
    USERS_POSTGRES_CONN_STR = "postgres://users:${var.users_password}@${var.pg_host}:${var.pg_port}/usersdb"
    RECEIPTS_POSTGRES_CONN_STR = "postgres://receipts:${var.receipts_password}@${var.pg_host}:${var.pg_port}/receiptsdb"
  }

}
