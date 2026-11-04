---
page_title: "polycrate_host Resource - polycrate"
subcategory: ""
description: |-
  A Polycrate Host (VM or bare metal).
---

# polycrate_host (Resource)

Maps 1:1 to `POST/GET/PATCH/DELETE /api/v1/hosts/`.

## Example Usage

```terraform
resource "polycrate_host" "web" {
  name         = "web-1"
  workspace_id = var.workspace_id
  product_id   = var.host_product_id
  kind         = "vm"
  role         = "generic"
}
```

## Schema

### Required

- `name` (String) Object name (slug). Forces replacement.
- `workspace_id` (String) Workspace UUID. Forces replacement.
- `product_id` (String) Product UUID with `kind=host`. Update PATCHes the Host.

### Optional

- `kind` (String) `vm` or `bare_metal`. Default `vm`. Forces replacement.
- `role` (String) Default `generic`.

### Read-Only

- `id` (String) Host UUID.
- `default_ipv4` (String) Public IPv4 once provisioned.
- `bootstrap_status` (String) API bootstrap status.

## Import

Import is supported using the Host UUID.

```shell
terraform import polycrate_host.web 11111111-1111-1111-1111-111111111111
```
