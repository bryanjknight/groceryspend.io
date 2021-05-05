resource "kubernetes_deployment" "predict" {
  metadata {
    name      = "predict"
    namespace = kubernetes_namespace.app.metadata.0.name
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "predict"
      }
    }
    template {
      metadata {
        labels = {
          app = "predict"
        }
      }
      spec {
        container {
          image = "groceryspend/predict:local"
          name  = "predict"
          port {
            container_port = 5000
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

resource "kubernetes_service" "predict" {
  metadata {
    name      = "predict-service"
    namespace = kubernetes_namespace.app.metadata.0.name
  }
  spec {
    selector = {
      app = kubernetes_deployment.predict.metadata.0.name
    }

    port {
      port        = 5000
      target_port = 5000
    }
  }
}

