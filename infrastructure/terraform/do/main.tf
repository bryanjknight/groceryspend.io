terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.4.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.0.1"
    }
  }
}

module "infra" {
  source               = "./infra"
  pvt_key              = var.pvt_key
  project_id           = data.digitalocean_project.groceryspend.id
  terraform_public_key = data.digitalocean_ssh_key.terraform.id
}

module "k8s-infra" {
  source = "./k8s-infra"
  vpc_id = module.infra.vpc_id
}

module "k8s-app" {
  source                        = "./k8s-app"
  cluster_name                  = module.k8s-infra.cluster_name
  cluster_id                    = module.k8s-infra.cluster_id
  cluster_endpoint              = module.k8s-infra.cluster_endpoint
  cluster_kubeconfig_raw_config = module.k8s-infra.cluster_kubeconfig_raw_config
  cluster_kubeconfig_token      = module.k8s-infra.cluster_kubeconfig_token
  cluster_kubeconfig_ca_cert    = module.k8s-infra.cluster_kubeconfig_ca_cert
  receipts_password             = module.infra.receipts_password
  users_password                = module.infra.users_password
  pg_host                       = module.infra.postgres_host
  pg_port                       = module.infra.postgres_port
}
