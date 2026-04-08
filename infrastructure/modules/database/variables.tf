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
  description = "Cloud SQL: PD_HDD or PD_SSD."
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

variable "private_service_networking_connection" {
  type        = any
  description = "google_service_networking_connection from the network module. Real resource reference for depends_on so destroy order is Cloud SQL → PSA, not parallel."
}

variable "db_name" {
  type = string
}

variable "db_user" {
  type = string
}
