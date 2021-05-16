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


First Time Setup
---
See [this document](https://github.com/bryanjknight/groceryspend.io/blob/main/infrastructure/docs/runbooks/first-time-deploy.md) for first time setup


Known Issues
---
* Existing bug where null resource and user creation don't work as expected. Requires deleting the users and re-running the plan and apply

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

docker tag web-portal:local groceryspend/web-portal:local
docker push groceryspend/web-portal:local
```

Notes
---
* IaC is not tied to a CICD pipeline as some minor changes result in the entire stack being torn down.
* Creating an ingress creates a load balancer in DO, resulting in another $10/mo for that


To Do
---
- [ ] Playbook for making incremental changes in UI, then capturing within Terraform config

References
---
- [Securing K8s](https://www.digitalocean.com/community/tutorials/recommended-steps-to-secure-a-digitalocean-kubernetes-cluster) to better secure it
