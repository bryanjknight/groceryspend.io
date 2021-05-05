terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.4.0"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.0.1"
    }
  }
}

resource "local_file" "kubeconfig" {
  depends_on = [var.cluster_id]
  count      = var.write_kubeconfig ? 1 : 0
  content    = var.cluster_kubeconfig_raw_config
  filename   = "${path.root}/kubeconfig"
}

provider "kubernetes" {
  host             = var.cluster_endpoint
  token            = var.cluster_kubeconfig_token
  cluster_ca_certificate = base64decode(
    var.cluster_kubeconfig_ca_cert
  )
}

provider "helm" {
  kubernetes {
    host  = var.cluster_endpoint
    token = var.cluster_kubeconfig_token
    cluster_ca_certificate = base64decode(
      var.cluster_kubeconfig_ca_cert
    )
  }
}