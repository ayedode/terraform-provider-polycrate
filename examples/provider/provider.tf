terraform {
  required_providers {
    polycrate = {
      source = "ayedo/polycrate"
    }
  }
}

provider "polycrate" {
  endpoint = "https://app.ayedo.cloud"
  # token    = var.polycrate_token  # or POLYCRATE_TOKEN
}
