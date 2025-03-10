# sopsage Provider

A Terraform provider for creating and managing SOPs Age Keys. See [Mozilla sops](https://github.com/mozilla/sops) for more info on SOPs Age.

!> WARNING This is a work in progress! Currently SOPs private keys are written to terraform state files in plain-text. I recommend that you encrypt your state before and after making additions to your project to prevent any damange of leaked Private SOPs Keys, Also recommend use a *secure* [state backend](https://www.terraform.io/docs/state/sensitive-data.html).

## Example Usage

```hcl
provider "sopsage" {}

# Creating a new SOPS Age key
resource "sopsage_key" "main" {
    provider = sopsage
}

# Output of Private Key Displayed (Returns value of `private_key = <sensitive>`)
output "private_key" {
  value = sopsage_key.main.private_key
  sensitive = true
}

# Output of Public Key Displayed
output "public_key" {
  value = sopsage_key.main.public_key
}
```