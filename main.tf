provider "azurerm" {
  features {}
}

module "rg_module" {
  source = "git::https://github.com/akashvarshney2023/terraform-module-rg.git?ref=1.0.0"
  resource_group_name = "test-resource-group"
  location = "westus2"
  tags = {
    environment = "production"
    owner       = "Akash"
  }
}
