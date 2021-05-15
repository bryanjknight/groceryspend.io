output "instance_id" {
  value = cloudamqp_instance.instance.id
}

output "rabbitmq_conn_str" {
  value     = cloudamqp_instance.instance.url
  sensitive = true
}
