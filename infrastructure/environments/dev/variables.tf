variable "project_id" {
  type = string
}

variable "github_repo_url" {
  type    = string
  default = "https://github.com/raphaelCamblong/iac-epitech.git"
}

variable "github_pat" {
  type      = string
  sensitive = true
}

variable "region" {
  type    = string
  default = "europe-west9"
}

variable "name_prefix" {
  type    = string
  default = "task-manager"
}

variable "subnet_cidr" {
  type    = string
  default = "10.10.0.0/20"
}

variable "pods_secondary_range_name" {
  type    = string
  default = "pods"
}

variable "pods_cidr" {
  type    = string
  default = "10.20.0.0/16"
}

variable "services_secondary_range_name" {
  type    = string
  default = "services"
}

variable "services_cidr" {
  type    = string
  default = "10.30.0.0/20"
}

variable "gke_release_channel" {
  type    = string
  default = "REGULAR"
}

variable "node_count" {
  type    = number
  default = 1
}

variable "node_machine_type" {
  type    = string
  default = "e2-standard-2"
}

variable "node_disk_size_gb" {
  type    = number
  default = 50
}

variable "db_instance_name" {
  type    = string
  default = "task-manager-db"
}

variable "db_tier" {
  type    = string
  default = "db-custom-1-3840"
}

variable "db_availability_type" {
  type    = string
  default = "ZONAL"
}

variable "db_disk_size_gb" {
  type    = number
  default = 20
}

variable "db_name" {
  type    = string
  default = "taskmanager"
}

variable "db_user" {
  type    = string
  default = "taskmanager"
}

variable "artifact_registry_name" {
  type    = string
  default = "task-manager"
}

variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "image_tag" {
  type = string
}

variable "ingress_host" {
  type = string
}

variable "replica_count" {
  type    = number
  default = 1
}

variable "hpa_min_replicas" {
  type    = number
  default = 1
}

variable "hpa_max_replicas" {
  type    = number
  default = 3
}
