provider "azurerm" {
  features {}
}



module "rg" {
  source  = "git::https://github.com/akashvarshney2023/terraform-module-rg.git?ref=1.0.0"
  generic = var.generic
  rg      = var.rg

}
