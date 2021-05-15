
module "cloudamqp-instance-infra" {
  source          = "./cloudamqp/instance"
  cloudamqp_token = var.cloudamqp_token
}

module "do-infra" {
  source            = "./do"
  pvt_key           = var.pvt_key
  do_token          = var.do_token
  rabbitmq_conn_str = module.cloudamqp-instance-infra.rabbitmq_conn_str
}

# Not available in free offering
# module "cloudamqp-fw-infra" {
#   source          = "./cloudamqp/fw"
#   cloudamqp_token = var.cloudamqp_token
# }