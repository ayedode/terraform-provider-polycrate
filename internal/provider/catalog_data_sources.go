// Copyright (c) Ayedo Cloud GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ayedode/terraform-provider-polycrate/internal/client"
)

func NewProductDataSource() datasource.DataSource { return &productDataSource{} }
func NewProviderDataSource() datasource.DataSource {
	return &catalogDataSource{typeName: "provider", fetch: fetchProvider}
}
func NewPopDataSource() datasource.DataSource {
	return &catalogDataSource{typeName: "pop", fetch: fetchPop}
}

type productDataSource struct {
	client *client.Client
}

type productModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Kind   types.String `tfsdk:"kind"`
	Active types.Bool   `tfsdk:"active"`
}

func (d *productDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

func (d *productDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up a Polycrate Product by UUID (`GET /api/v1/products/{id}/`). Use Host SKUs (`kind=host`).",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Required: true, MarkdownDescription: "Product UUID."},
			"name":   schema.StringAttribute{Computed: true, MarkdownDescription: "Product name."},
			"kind":   schema.StringAttribute{Computed: true, MarkdownDescription: "Product kind (`host` for Compute Cloud SKUs)."},
			"active": schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the product is active."},
		},
	}
}

func (d *productDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *productDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data productModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	p, err := d.client.GetProduct(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read product", err.Error())
		return
	}
	data.Name = types.StringValue(p.Name)
	data.Kind = types.StringValue(p.Kind)
	data.Active = types.BoolValue(p.Active)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type catalogDataSource struct {
	client   *client.Client
	typeName string
	fetch    func(context.Context, *client.Client, string) (*client.CatalogObject, error)
}

type catalogModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	SalesReady  types.Bool   `tfsdk:"sales_ready"`
	SCSIaaS     types.String `tfsdk:"scs_iaas_status"`
	ISO27001    types.String `tfsdk:"iso27001_status"`
}

func fetchProvider(ctx context.Context, c *client.Client, id string) (*client.CatalogObject, error) {
	return c.GetProvider(ctx, id)
}

func fetchPop(ctx context.Context, c *client.Client, id string) (*client.CatalogObject, error) {
	return c.GetPop(ctx, id)
}

func (d *catalogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.typeName
}

func (d *catalogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Look up a Polycrate %s by UUID.", d.typeName),
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Required: true, MarkdownDescription: "Object UUID."},
			"name":            schema.StringAttribute{Computed: true, MarkdownDescription: "Object name."},
			"display_name":    schema.StringAttribute{Computed: true, MarkdownDescription: "Display name."},
			"sales_ready":     schema.BoolAttribute{Computed: true, MarkdownDescription: "Sales-ready catalog flag."},
			"scs_iaas_status": schema.StringAttribute{Computed: true, MarkdownDescription: "SCS IaaS status (`compatible`, `listed_not_certified`, `none`)."},
			"iso27001_status": schema.StringAttribute{Computed: true, MarkdownDescription: "ISO 27001 catalog status."},
		},
	}
}

func (d *catalogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *catalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data catalogModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	o, err := d.fetch(ctx, d.client, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read "+d.typeName, err.Error())
		return
	}
	data.Name = types.StringValue(o.Name)
	data.DisplayName = types.StringValue(o.DisplayName)
	data.SalesReady = types.BoolValue(o.SalesReady)
	data.SCSIaaS = types.StringValue(o.SCSIaaS)
	data.ISO27001 = types.StringValue(o.ISO27001)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
