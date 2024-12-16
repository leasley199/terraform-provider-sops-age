package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.mozilla.org/sops/decrypt"
)

func dataSourceSopsAgeFile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSopsAgeFileRead,
		Schema: map[string]*schema.Schema{
			"file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the SOPS-encrypted file",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Decrypted content of the file",
			},
		},
	}
}

func dataSourceSopsAgeFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	filePath := d.Get("file_path").(string)
	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read file: %w", err))
	}

	decryptedData, err := decrypt.Data(encryptedData, "yaml")
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to decrypt file: %w", err))
	}

	d.Set("content", string(decryptedData))
	d.SetId(filePath) // Use the file path as a unique ID

	return diags
}
