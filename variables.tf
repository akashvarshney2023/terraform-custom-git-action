variable "resource_group_name" {
  type        = string
  description = "name of the resource group"
}

variable "rg_location" {
  type        = string
  description = "The Azure region where the resources will be created"
}

variable "vnet_name" {
  description = "The name of the VNet to create."
  type        = string
}

variable "address_space" {
  description = "The CIDR block for the VNet address space."
  type        = list(string)
}

variable "subnets" {
  description = "A list of subnet configurations."
  type        = list(object({
    name           = string
    cidr_block     = string
    service_endpoints = list(string)
  }))
}
