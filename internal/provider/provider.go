// Copyright (c) Ayedo Cloud GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ayedode/terraform-provider-polycrate/internal/client"
)

var _ provider.Provider = &PolycrateProvider{}

type PolycrateProvider struct {
	version string
}

type PolycrateProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p *PolycrateProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "polycrate"
	resp.Version = p.version
}

func (p *PolycrateProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Polycrate Hosts (ayedo Compute Cloud) through the existing Host REST API. " +
			"Works with Terraform and OpenTofu. There is no second VM resource type.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Polycrate API base URL, without a trailing path. Example: `https://app.ayedo.cloud`. " +
					"May also be set with `POLYCRATE_ENDPOINT`.",
				Optional: true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "API token (User, Organization, or System key). Sent as `Authorization: Bearer`. " +
					"May also be set with `POLYCRATE_TOKEN`.",
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *PolycrateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PolycrateProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("POLYCRATE_ENDPOINT")
	token := os.Getenv("POLYCRATE_TOKEN")
	if !data.Endpoint.IsNull() && !data.Endpoint.IsUnknown() {
		endpoint = data.Endpoint.ValueString()
	}
	if !data.Token.IsNull() && !data.Token.IsUnknown() {
		token = data.Token.ValueString()
	}
	if endpoint == "" {
		resp.Diagnostics.AddError("Missing endpoint", "Set provider.endpoint or POLYCRATE_ENDPOINT.")
		return
	}
	if token == "" {
		resp.Diagnostics.AddError("Missing token", "Set provider.token or POLYCRATE_TOKEN.")
		return
	}

	c := client.New(endpoint, token)
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *PolycrateProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHostResource,
	}
}

func (p *PolycrateProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProductDataSource,
		NewProviderDataSource,
		NewPopDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PolycrateProvider{version: version}
	}
}
