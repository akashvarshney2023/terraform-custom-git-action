provider "azurerm" {
  features {

  }

}

module "resource_rg" {
  source                         = "./modules/rg"
  module_resource_group_name     = var.rg_name
  module_resource_group_location = var.rg_location
}


module "vnet" {
  source              = "./modules/vnet"
  location            = var.rg_location
  resource_group_name = var.resource_group_name
  vnet_name           = var.vnet_name
  address_space       = var.address_space
  subnet              = var.subnets
}
