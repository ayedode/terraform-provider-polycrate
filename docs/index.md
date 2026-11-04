---
page_title: "Polycrate Provider"
description: |-
  Manage Polycrate Hosts through the Polycrate REST API.
---

# Polycrate Provider

Terraform and OpenTofu provider for [Polycrate](https://docs.ayedo.de/polycrate/). Resource `polycrate_host` maps to `/api/v1/hosts/`. There is no `hetzner_*` or `openstack_*` resource.

## Example Usage

```terraform
provider "polycrate" {
  endpoint = "https://app.ayedo.cloud"
}
```

## Schema

### Optional

- `endpoint` (String) API base URL. Environment: `POLYCRATE_ENDPOINT`.
- `token` (String, Sensitive) Bearer token. Environment: `POLYCRATE_TOKEN`.
