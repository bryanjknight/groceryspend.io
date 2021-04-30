resource "digitalocean_droplet" "bastion" {
  image = "ubuntu-20-04-x64"
  name = "bastion"
  region = "${var.region}"
  size = "${var.droplet_image}"
  
  vpc_uuid = digitalocean_vpc.staging.id
  private_networking = true

  ssh_keys = [
    data.digitalocean_ssh_key.terraform.id
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
      # install nginx
      "sudo apt-get update",
    ]
  }
}

resource "digitalocean_project_resources" "bastion_groceryspend" {
  project = data.digitalocean_project.groceryspend.id
  resources = [
    digitalocean_droplet.bastion.urn
  ]
}