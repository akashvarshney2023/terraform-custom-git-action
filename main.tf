provider "azurerm" {
  features {}
}



module "rg" {
  source  = "git::https://github.com/akashvarshney2023/terraform-module-rg.git?ref=1.0.0"
  generic = var.generic
  rg      = var.rg

}

module "resource_group1_test" {
  source  = "git::https://github.com/akashvarshney2023/terraform-module-rg.git?ref=1.0.0"
  generic = var.generic
  rg      = var.rg

}


module "resource_group2_test" {
  source  = "git::https://github.com/akashvarshney2023/terraform-module-rg.git?ref=1.0.0"
  generic = var.generic
  rg      = var.rg

}