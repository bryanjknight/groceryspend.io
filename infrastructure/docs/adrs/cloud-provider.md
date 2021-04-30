Cloud Provider
===

## Status
Approved 2021-04-30 

## Context
The goal of this ADR is to pick a cloud provider for a managed k8s cluster, priortizing costs over functionality and availablity

## Options

We are using the following criteria for running k8s
* Master should have 2 vCPUs, 2 GB RAM at least ([link here](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/))
* Workers should have the same

| Cloud Provider | Master Node | Spot Price for workers | Total Cost unmanaged | Total Cost managed |
| ---------------| ----------- | ---------------------- | -------------------- | ------------------ |
| [AWS](https://aws.amazon.com/ec2/spot/pricing/) | t4g.small	@ 0.0168/hr => 12.26/mo | a1.medium @ 0.0049/hr 3.57/mo | $15.83/mo | $88.83/mo |
| [Azure](https://azure.microsoft.com/en-us/pricing/details/virtual-machines/linux/) | B2S	@ 30.42/mo | B2S @ 12.1472/mo | $42.56/mo | $42.56/mo  (AKS is free)|
| [GCP](https://cloud.google.com/products/calculator#id=a9b713ee-a39a-4e13-af6b-c38398a213ec) | e2-standard-2 @ 48.92/mo | 14.68/mo | 78.19/mo | 78.19/mo (one zonal gke free) |
| DO | $15/mo | $15/mo | $30/mo | $20/mo |
| Linode | $10/mo | $10/mo| $20/mo | $20/mo |


[A review of GCP, DO, and Linode](https://atodorov.me/2020/06/14/comparing-kubernetes-managed-services-across-digital-ocean-scaleway-ovhcloud-and-linode/#conclusions) showed that DO had more features for the same price as Linode.

## Decision
We'll use DigitialOcean for our K8S cluster (~$20/mo), with Docker Hub (free for public repos) as I'm not concerned about the potential of people downloading the image. We can futher cut costs by reducing the number of nodes to 1, removing load balancing.


## Consequences
DigitalOcean has some features missing and we'll need to spike to figure out how to manage this all properly

## Compliance
N/A

## Notes
N/A