---
page_title: "polycrate_provider Data Source - polycrate"
description: |-
  Look up a Polycrate Provider by UUID.
---

# polycrate_provider (Data Source)

`GET /api/v1/providers/{id}/`.

## Schema

### Required

- `id` (String)

### Read-Only

- `name` (String)
- `display_name` (String)
- `sales_ready` (Boolean)
- `scs_iaas_status` (String)
- `iso27001_status` (String)
