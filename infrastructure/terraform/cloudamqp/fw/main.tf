terraform {
  required_providers {
    cloudamqp = {
      source  = "cloudamqp/cloudamqp"
      version = ">= 1.9.1"
    }
  }
}

# Configure the CloudAMQP Provider
provider "cloudamqp" {
  apikey        = var.cloudamqp_token
}

# 
# NOTE: this is only available on dedicated instances
#
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = var.instance_id

  rules {
    ip          = "192.168.0.0/24"
    ports       = [4567, 4568]
    services    = ["AMQP","AMQPS"]
  }

  rules {
    ip          = "10.56.72.0/24"
    ports       = []
    services    = ["AMQP","AMQPS"]
  }
}
