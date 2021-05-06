# first install nginx ingress controller, which will create an ingress service
resource "helm_release" "nginx_ingress" {
  name       = "nginx-ingress-controller"
  namespace  = kubernetes_namespace.app.metadata.0.name

  repository = "https://charts.bitnami.com/bitnami"
  chart      = "nginx-ingress-controller"

  # we want a load balancers
  set {
    name  = "service.type"
    value = "LoadBalancer"
  }

  # set the annotations
  set {
    name = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/do-loadbalancer-protocol"
    value = "http"
  }
  set {
    name = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/do-loadbalancer-size-slug"
    value = "lb-small"
  }
    set {
    name = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/do-loadbalancer-name"
    value = "main-lb"
  }

}
