// Copyright (c) Ayedo Cloud GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostResource(t *testing.T) {
	testAccPreCheck(t)

	workspace := os.Getenv("POLYCRATE_WORKSPACE_ID")
	product := os.Getenv("POLYCRATE_PRODUCT_ID")
	name := fmt.Sprintf("tf-acc-%d", os.Getpid())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "polycrate" {}

resource "polycrate_host" "test" {
  name         = %q
  workspace_id = %q
  product_id   = %q
  kind         = "vm"
  role         = "generic"
}
`, name, workspace, product),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("polycrate_host.test", "name", name),
					resource.TestCheckResourceAttr("polycrate_host.test", "role", "generic"),
					resource.TestCheckResourceAttrSet("polycrate_host.test", "id"),
				),
			},
			{
				ResourceName:      "polycrate_host.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
