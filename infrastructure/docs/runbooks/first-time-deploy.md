First Time Deployment Checklist
===

- [ ] Setup Auth0
  - [ ] Enusre domains are valid in callback
- [ ] Setup DO account manually
  - [ ] Project
  - [ ] SSH Key
  - [ ] Domain
- [ ] Deploy infrastructure via terraform
  - [ ] `cd ./infrastructure/terraform/do && terraform init`
  - [ ] `terraform plan -out infra.out`
  - [ ] `terraform apply infra.out`
- [ ] Restart deploys from CICD (this will deploy the deployments to k8s)
- [ ] Deploy k8s features
  - [ ] `./infrastructure/k8s/init_deploy.sh`
- [ ] Tie load balancer IP to DNS records via A records
  - [ ] `@`
  - [ ] `api` 
  - [ ] `www`
  - [ ] `portal`
- [ ] Create SSL cert for load balancers
  - [ ] Determine https port endpoint for ingress by going to k8s console -> service -> services -> nginx-ingress-controller -> triple dots. Find http (**NOT** HTTPS) nodePort
  - [ ] Networking -> Load Balancers -> <random id> -> Settings -> Forwarding Rules. HTTP2 443 (select new cert, then select domain and use for all subdomains), 