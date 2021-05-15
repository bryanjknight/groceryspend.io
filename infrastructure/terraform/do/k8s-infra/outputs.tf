output "cluster_id" {
  value = digitalocean_kubernetes_cluster.groceryspend.id
}

output "cluster_name" {
  value = digitalocean_kubernetes_cluster.groceryspend.name
}

output "cluster_endpoint" {
  value     = digitalocean_kubernetes_cluster.groceryspend.endpoint
  sensitive = true
}

output "cluster_kubeconfig_raw_config" {
  value     = digitalocean_kubernetes_cluster.groceryspend.kube_config[0].raw_config
  sensitive = true
}

output "cluster_kubeconfig_token" {
  value     = digitalocean_kubernetes_cluster.groceryspend.kube_config[0].token
  sensitive = true
}

output "cluster_kubeconfig_ca_cert" {
  value     = digitalocean_kubernetes_cluster.groceryspend.kube_config[0].cluster_ca_certificate
  sensitive = true
}

output "cluster_tags" {
  value = digitalocean_kubernetes_cluster.groceryspend.node_pool.0.tags
}
