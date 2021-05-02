Infrastructure
===

Repo for all terraform logic to spin up/down environment

Manually Created Resources
---

### AWS
* Domain name registration

### DigitalOcean
* Project
* SSH Key
* DNS records
* K8S
  * Maintenance window
  * Burst control

Tools
---
- doctl: `brew install doctl`
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

Docker deploy process w/ DigitalOcean
---
1. Install doctl and docker
2. Authenticate with doctl with API token
3. doctl registry login
4. Build images
5. Run the following:

```
docker tag predict:local registry.digitalocean.com/groceryspend/predict:local
docker push registry.digitalocean.com/groceryspend/predict:local

docker tag server:local registry.digitalocean.com/groceryspend/server:local
docker push registry.digitalocean.com/groceryspend/server:local
```

Docker deploy process w/ Docker
---
1. Install docker
2. Authenticate with docker
4. Build images
5. Run the following:

```
docker tag predict:local groceryspend/predict:local
docker push groceryspend/predict:local

docker tag server:local groceryspend/server:local
docker push groceryspend/server:local
```

Notes
---
* IaC is not tied to a CICD pipeline as some minor changes result in the entire stack being torn down.
* Creating an ingress creates a load balancer in DO, resulting in another $10/mo for that

To Do
---
- [ ] Playbook for making incremental changes in UI, then capturing within Terraform config
