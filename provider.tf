terraform {
  backend "azurerm" {
    resource_group_name  = "NetworkWatcherRG"
    storage_account_name = "terraform24042023"
    container_name       = "terraform"
    key                  = "terraform.tfstate"
  }
}