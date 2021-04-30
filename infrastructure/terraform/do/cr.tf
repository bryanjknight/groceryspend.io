# Create a new container registry
resource "digitalocean_container_registry" "groceryspend_cr" {
  name                   = "groceryspend"
  subscription_tier_slug = "starter"
}