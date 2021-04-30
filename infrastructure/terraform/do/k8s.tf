resource "digitalocean_kubernetes_cluster" "groceryspend_test" {
  name   = "groceryspend-test"
  region = "${var.region}"
  # Grab the latest version slug from `doctl kubernetes options versions`
  version = "1.20.2-do.0"

  vpc_uuid = digitalocean_vpc.staging.id

  node_pool {
    name       = "worker-pool"
    size       = "${var.k8s_worker_image}"
    node_count = 1

    tags = [ "k8s" ]
    taint {
      key    = "workloadKind"
      value  = "database"
      effect = "NoSchedule"
    }
  }
}
