package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = (*carstoreProvider)(nil)

func New() func() provider.Provider {
	return func() provider.Provider {
		return &carstoreProvider{}
	}
}

type carstoreProvider struct {
	baseURL string
	client  *http.Client
}

type carstoreProviderModel struct {
	BaseURL types.String `tfsdk:"base_url"`
}

func (p *carstoreProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *carstoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config carstoreProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configure the provider with a default timeout client
	p.client = &http.Client{
		Timeout: 10 * time.Second,
	}
	p.baseURL = config.BaseURL.ValueString()
}

func (p *carstoreProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "carstore"
}

func (p *carstoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *carstoreProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return NewCarResource(p.baseURL, p.client)
		},
	}
}
