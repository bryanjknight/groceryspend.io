data "digitalocean_loadbalancer" "lb" {
  name = "groceryspend.io"
  depends_on = [
    helm_release.nginx_ingress,
    kubernetes_ingress.groceryspend_ingress,
  ]
}

output "load_balancer_ip" {
  value = data.digitalocean_loadbalancer.lb.ip
}
