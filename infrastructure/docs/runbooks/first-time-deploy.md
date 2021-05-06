First Time Deployment
===

- [ ] Setup DO account manually
  - [ ] Project
  - [ ] SSH Key
  - [ ] Domain
- [ ] Deploy infrastructure via terraform
  - [ ] `cd ./infrastructure/terraform/do && terraform init`
  - [ ] `terraform plan -out infra.out`
  - [ ] `terraform apply infra.out`
- [ ] Deploy k8s features
  - [ ] `./infrastructure/k8s/init_deploy.sh`
- [ ] Tie load balancer IP to DNS records via A records
  - [ ] `@`
  - [ ] `api` 
  - [ ] `www`
  - [ ] `portal`