resource "kubernetes_config_map" "server-config" {
  metadata {
    name      = "server-config"
    namespace = kubernetes_namespace.app.metadata.0.name
  }

  data = {
    AUTH_PROVIDER                 = "AUTH0"
    AUTH_ALLOW_ORIGINS            = "http://localhost:3000 chrome-extension://gpmoghmaibomfddfbofkionknjjeoaef"
    AUTH_ALLOW_METHODS            = "GET PUT POST DELETE PATCH OPTIONS"
    AUTH_ALLOW_HEADERS            = "* Authorization"
    AUTH_EXPOSE_HEADERS           = "*"
    AUTH_ALLOW_CREDENTIALS        = "true"
    AUTH_ALLOW_BROWSER_EXTENSIONS = "true"
    AUTH_MAX_AGE                  = "12h"
    CATEGORIZE_PATH               = "/categorize"
  }
}
