Infrastructure
===

Repo for all terraform logic to spin up/down environment

Manually Created Resources
---

### AWS
* Domain name

### DigitalOcean
* Project
* SSH Key
* DNS records

Tools
---
- Terraform: `brew tap hashicorp/tap && brew install hashicorp/tap/terraform`
- kubectl `brew install kubernetes-cli`
- helm `brew install helm`


Setup
---
1. Get the k8s access config, copy it into `$HOME/.kube/config`


Things to automate:
* `helm repo add nginx-stable https://helm.nginx.com/stable`
* `helm repo update`
* `helm install nginx-ingress-v1 nginx-stable/nginx-ingress`
* `kubectl apply -f ./namespaces.yml`
* `kubectl apply -f ./groceryspend/hello-world.yml`
* `kubectl apply -f ./groceryspend/ingress.yml`
* `helm repo add honeycomb https://honeycombio.github.io/helm-charts`
* `helm install honeycomb honeycomb/honeycomb --set honeycomb.apiKey=API_KEY`

Notes
---
* IaC is not tied to a CICD pipeline as some minor changes result in the entire stack being torn down.
* Creating an ingress creates a load balancer in DO, resulting in another $10/mo for that

To Do
---
- [ ] Playbook for making incremental changes in UI, then capturing within Terraform config
