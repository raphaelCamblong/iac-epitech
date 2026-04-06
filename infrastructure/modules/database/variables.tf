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

variable "db_disk_type" {
  type        = string
  description = "Use 'PD_STANDARD' for HDD (low footprint/no quota issues) or 'PD_SSD' for high performance."
  default     = "PD_HDD"
}

variable "db_disk_autoresize" {
  type    = bool
  default = true
}

variable "db_disk_autoresize_limit_gb" {
  type        = number
  description = "Cap for autoresize; 0 = provider default (no cap)."
  default     = 20
}

variable "network_id" {
  type = string
}

variable "private_vpc_connection" {
  type        = string
  description = "Dependency on private VPC connection"
}

variable "db_name" {
  type = string
}

variable "db_user" {
  type = string
}
