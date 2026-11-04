resource "polycrate_host" "web" {
  name         = "web-1"
  workspace_id = var.workspace_id
  product_id   = var.host_product_id
  kind         = "vm"
  role         = "generic"
}
