package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/leasley199/terraform-provider-sops-age/sopsage"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderAddr: "registry.terraform.io/leasley199/sops-age",
		ProviderFunc: func() *schema.Provider {
			return sopsage.Provider()
		},
	})
}
