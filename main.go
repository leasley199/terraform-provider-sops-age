package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{

		ProviderFunc: func() *schema.Provider {
			return &schema.Provider{
				Schema: map[string]*schema.Schema{
					"passphrase": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
				ResourcesMap: map[string]*schema.Resource{
					"sops_age_key": resourceSopsAgeKey(),
				},
				DataSourcesMap: map[string]*schema.Resource{
					"sops_age_file": dataSourceSopsAgeFile(),
				},
			}
		},
	})
}
