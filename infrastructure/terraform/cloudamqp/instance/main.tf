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

# Create a new cloudamqp instance
resource "cloudamqp_instance" "instance" {
  name          = "terraform-cloudamqp-instance"
  plan          = "lemur"
  region        = "amazon-web-services::us-east-1"
  nodes         = 1
  tags          = [ "terraform" ]
  rmq_version   = "3.8.3"
}

