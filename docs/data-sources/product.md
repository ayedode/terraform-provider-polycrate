---
page_title: "polycrate_product Data Source - polycrate"
description: |-
  Look up a Polycrate Product by UUID.
---

# polycrate_product (Data Source)

`GET /api/v1/products/{id}/`.

## Example Usage

```terraform
data "polycrate_product" "sku" {
  id = var.host_product_id
}
```

## Schema

### Required

- `id` (String) Product UUID.

### Read-Only

- `name` (String)
- `kind` (String)
- `active` (Boolean)
