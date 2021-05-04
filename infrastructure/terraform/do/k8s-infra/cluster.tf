resource "digitalocean_kubernetes_cluster" "groceryspend" {
  name   = "groceryspend-${var.namespace}"
  region = "${var.region}"
  # Grab the latest version slug from `doctl kubernetes options versions`
  version = "${var.k8s_version}"

  vpc_uuid = "${var.vpc_id}"

  node_pool {
    name       = "worker-pool"
    size       = "${var.k8s_worker_image}"
    node_count = "${var.k8s_node_count}"

    tags = [ "k8s" ]
    taint {
      key    = "workloadKind"
      value  = "database"
      effect = "NoSchedule"
    }
  }
}

