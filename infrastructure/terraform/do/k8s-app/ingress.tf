# first install nginx ingress controller
resource "helm_release" "nginx_ingress" {
  name       = "nginx-ingress-controller"
  namespace  = kubernetes_namespace.app.metadata.0.name

  repository = "https://charts.bitnami.com/bitnami"
  chart      = "nginx-ingress-controller"

  set {
    name  = "service.type"
    value = "LoadBalancer"
  }
}

# then create the ingress itself
resource "kubernetes_ingress" "groceryspend_ingress" {
  wait_for_load_balancer = true
  metadata {
    name = "ingress"
    namespace  = kubernetes_namespace.app.metadata.0.name
    annotations = {
      "kubernetes.io/ingress.class" = "nginx"
    }
  }

  spec {
    rule {
      host = "www.groceryspend.io"
      http {
        path {
          backend {
            service_name = kubernetes_service.server.metadata.0.name
            service_port = 8080
          }

          path = "/"
        }
      }
    }
  }
}