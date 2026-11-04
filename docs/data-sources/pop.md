---
page_title: "polycrate_pop Data Source - polycrate"
description: |-
  Look up a Polycrate PoP by UUID.
---

# polycrate_pop (Data Source)

`GET /api/v1/pops/{id}/`.

## Schema

### Required

- `id` (String)

### Read-Only

- `name` (String)
- `display_name` (String)
- `sales_ready` (Boolean)
- `scs_iaas_status` (String)
- `iso27001_status` (String)
