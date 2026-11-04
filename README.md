# Terraform / OpenTofu Provider for Polycrate

Based on the official [HashiCorp Terraform Plugin Framework scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding-framework).

This provider manages **Polycrate Hosts** (`polycrate_host`) against `POST/GET/PATCH/DELETE /api/v1/hosts/`. It does not talk to Hetzner, OpenStack, or other IaaS APIs. Provider differences stay in polycrate-api drivers.

Compatible with Terraform >= 1.0 and OpenTofu (protocol 6.0).

GitHub publishing, GPG signing, and Terraform Registry / OpenTofu Registry listing are out of band (see `.github/workflows/release.yml`).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0 or [OpenTofu](https://opentofu.org/)
- [Go](https://go.dev/doc/install) >= 1.24 (to build)

## Build

```shell
go install
```

Dev override (OpenTofu / Terraform CLI config):

```hcl
provider_installation {
  dev_overrides {
    "ayedo/polycrate" = "/path/to/go/bin"
  }
  direct {}
}
```

## Configuration

```hcl
provider "polycrate" {
  endpoint = "https://app.ayedo.cloud"
  token    = var.polycrate_token
}
```

Environment: `POLYCRATE_ENDPOINT`, `POLYCRATE_TOKEN`. Token is a User, Organization, or System API key (`Authorization: Bearer`), same as [polycrate-sdk](https://github.com/ayedode/polycrate-sdk).

## Resources and data sources

| Type | API |
|------|-----|
| `polycrate_host` | `/api/v1/hosts/` |
| `polycrate_product` (data) | `/api/v1/products/{id}/` |
| `polycrate_provider` (data) | `/api/v1/providers/{id}/` |
| `polycrate_pop` (data) | `/api/v1/pops/{id}/` |

Lifecycle actions (`start`, `shutdown`, `reboot`, `reinstall`) remain API/MCP actions, not Terraform CRUD.

## Tests

```shell
go test ./...
# Acceptance (needs a live API):
POLYCRATE_ENDPOINT=https://app.ayedo.cloud POLYCRATE_TOKEN=... \
  POLYCRATE_WORKSPACE_ID=... POLYCRATE_PRODUCT_ID=... make testacc
```
