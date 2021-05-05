resource "kubernetes_deployment" "server" {
  metadata {
    name      = "server"
    namespace = kubernetes_namespace.app.metadata.0.name
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "server"
      }
    }
    template {
      metadata {
        labels = {
          app = "server"
        }
      }
      spec {
        container {
          image = "groceryspend/server:local"
          name  = "server"
          port {
            container_port = 8080
          }
          env_from {
            secret_ref {
              name = "server-auth"
            }
            config_map_ref {
              name = "server-config"
            }
          }
          resources {
            limits = {
              memory = "512M"
              cpu    = "1"
            }
            requests = {
              memory = "256M"
              cpu    = "50m"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "server" {
  metadata {
    name      = "server-service"
    namespace = kubernetes_namespace.app.metadata.0.name
  }
  spec {
    selector = {
      app = kubernetes_deployment.server.metadata.0.name
    }

    port {
      port        = 8080
      target_port = 8080
    }
  }
}
