variable "db_instance_name" {
  type = string
}

variable "region" {
  type = string
}

variable "db_tier" {
  type = string
}

variable "db_availability_type" {
  type = string
}

variable "db_disk_size_gb" {
  type = number
}

variable "network_id" {
  type = string
}

variable "private_vpc_connection" {
  type = string
  description = "Dependency on private VPC connection"
}

variable "db_name" {
  type = string
}

variable "db_user" {
  type = string
}
