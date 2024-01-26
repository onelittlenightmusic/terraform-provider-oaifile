// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure OAIFileProvider satisfies various provider interfaces.
var _ provider.Provider = &OAIFileProvider{}

// OAIFileProvider defines the provider implementation.
type OAIFileProvider struct {
	version	string
}

// OAIFileProviderModel describes the provider data model.
type OAIFileProviderModel struct {
	Host     types.String `tfsdk:"host"`
}

func (p *OAIFileProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "oaifile"
	resp.Version = p.version
}

func (p *OAIFileProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *OAIFileProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OAIFileProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Host.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
					path.Root("host"),
					"Unknown OAIFile API Host",
					"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API host. "+
							"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
			)
	}	// Example client configuration for data sources and resources
	
	if resp.Diagnostics.HasError() {
		return
	}
	client := &OAIFileClient{
		Host: data.Host.ValueString(),
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OAIFileProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOAIFileResource,
	}
}

func (p *OAIFileProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OAIFileProvider{
			version: version,
		}
	}
}
