// Copyright (c) Ayedo Cloud GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ayedode/terraform-provider-polycrate/internal/client"
)

var _ resource.Resource = &HostResource{}
var _ resource.ResourceWithImportState = &HostResource{}

func NewHostResource() resource.Resource {
	return &HostResource{}
}

type HostResource struct {
	client *client.Client
}

type HostResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	WorkspaceID     types.String `tfsdk:"workspace_id"`
	ProductID       types.String `tfsdk:"product_id"`
	Kind            types.String `tfsdk:"kind"`
	Role            types.String `tfsdk:"role"`
	DefaultIPv4     types.String `tfsdk:"default_ipv4"`
	BootstrapStatus types.String `tfsdk:"bootstrap_status"`
}

func (r *HostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *HostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A Polycrate Host (VM or bare metal). Maps 1:1 to `POST/GET/PATCH/DELETE /api/v1/hosts/`. " +
			"This is not a provider-native instance resource (no `hetzner_*` / `openstack_*`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Host UUID from the API. Terraform state key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Object name (slug).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Workspace UUID. Required on the Host broker path.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"product_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Product UUID with `kind=host` (SKU). Changing this PATCHes the Host (resize when the driver supports it).",
			},
			"kind": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("vm"),
				MarkdownDescription: "`vm` or `bare_metal`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("generic"),
				MarkdownDescription: "Host role. Default `generic` (raw VM). Kubernetes workers are created by WorkerPool, not this resource.",
			},
			"default_ipv4": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Public IPv4 assigned by the provider. Empty until provisioned.",
			},
			"bootstrap_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Bootstrap status from the API (`pending`, `provisioning`, `ready`, …).",
			},
		},
	}
}

func (r *HostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *HostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	h, err := r.client.CreateHost(ctx, map[string]any{
		"name":         data.Name.ValueString(),
		"workspace_id": data.WorkspaceID.ValueString(),
		"product_id":   data.ProductID.ValueString(),
		"kind":         data.Kind.ValueString(),
		"role":         data.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create host", err.Error())
		return
	}
	flattenHost(h, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	h, err := r.client.GetHost(ctx, data.ID.ValueString())
	if err != nil {
		if _, ok := err.(*client.NotFoundError); ok {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read host", err.Error())
		return
	}
	flattenHost(h, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data HostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	h, err := r.client.UpdateHost(ctx, data.ID.ValueString(), map[string]any{
		"product_id": data.ProductID.ValueString(),
		"role":       data.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to update host", err.Error())
		return
	}
	flattenHost(h, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteHost(ctx, data.ID.ValueString()); err != nil {
		if _, ok := err.(*client.NotFoundError); ok {
			return
		}
		resp.Diagnostics.AddError("Unable to delete host", err.Error())
	}
}

func (r *HostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func flattenHost(h *client.Host, data *HostResourceModel) {
	data.ID = types.StringValue(h.ID)
	data.Name = types.StringValue(h.Name)
	if h.Kind != "" {
		data.Kind = types.StringValue(h.Kind)
	}
	if h.Role != "" {
		data.Role = types.StringValue(h.Role)
	}
	if h.ProductID != "" {
		data.ProductID = types.StringValue(h.ProductID)
	}
	if h.WorkspaceID != "" {
		data.WorkspaceID = types.StringValue(h.WorkspaceID)
	}
	data.BootstrapStatus = types.StringValue(h.BootstrapStatus)
	if h.DefaultIPv4 != nil {
		data.DefaultIPv4 = types.StringValue(*h.DefaultIPv4)
	} else {
		data.DefaultIPv4 = types.StringValue("")
	}
}
