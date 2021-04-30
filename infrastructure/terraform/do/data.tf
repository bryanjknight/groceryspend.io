# Things that are alrady in DigitialOcean that we need to reference
data "digitalocean_ssh_key" "terraform" {
  # TODO: create a new key for managing the account
  name = "Digital Ocean - Home Laptop"
}

data "digitalocean_project" "groceryspend" {
  name = "groceryspend"
}
