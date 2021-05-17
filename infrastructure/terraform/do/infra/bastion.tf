resource "digitalocean_droplet" "bastion" {
  image = "ubuntu-20-04-x64"
  name = "bastion"
  region = "${var.region}"
  size = "${var.droplet_image}"
  
  vpc_uuid = digitalocean_vpc.groceryspend-vpc.id
  private_networking = true

  ssh_keys = [
    var.terraform_public_key
  ]

  tags = [ "bastion" ]
  connection {
    host = self.ipv4_address
    user = "root"
    type = "ssh"
    private_key = file(var.pvt_key)
    timeout = "2m"
  }
  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      "sudo apt-get update",
    ]
  }
}

resource "digitalocean_project_resources" "bastion_groceryspend" {
  project = var.project_id
  resources = [
    digitalocean_droplet.bastion.urn
  ]
}