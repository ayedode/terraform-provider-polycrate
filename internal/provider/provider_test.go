// Copyright (c) Ayedo Cloud GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"polycrate": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("POLYCRATE_ENDPOINT") == "" || os.Getenv("POLYCRATE_TOKEN") == "" {
		t.Skip("POLYCRATE_ENDPOINT and POLYCRATE_TOKEN must be set for acceptance tests")
	}
	if os.Getenv("POLYCRATE_WORKSPACE_ID") == "" || os.Getenv("POLYCRATE_PRODUCT_ID") == "" {
		t.Skip("POLYCRATE_WORKSPACE_ID and POLYCRATE_PRODUCT_ID must be set for acceptance tests")
	}
}
