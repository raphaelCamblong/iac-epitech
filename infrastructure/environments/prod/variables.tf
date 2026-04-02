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
  default = "10.0.0.0/20"
}

variable "pods_secondary_range_name" {
  type    = string
  default = "pods"
}

variable "pods_cidr" {
  type    = string
  default = "10.1.0.0/16"
}

variable "services_secondary_range_name" {
  type    = string
  default = "services"
}

variable "services_cidr" {
  type    = string
  default = "10.2.0.0/20"
}

variable "gke_release_channel" {
  type    = string
  default = "STABLE"
}

variable "node_count" {
  type    = number
  default = 3
}

variable "node_machine_type" {
  type    = string
  default = "e2-standard-4"
}

variable "node_disk_size_gb" {
  type    = number
  default = 100
}

variable "db_instance_name" {
  type    = string
  default = "task-manager-db"
}

variable "db_tier" {
  type    = string
  default = "db-custom-2-7680"
}

variable "db_availability_type" {
  type    = string
  default = "REGIONAL"
}

variable "db_disk_size_gb" {
  type    = number
  default = 100
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
  default = 3
}

variable "hpa_min_replicas" {
  type    = number
  default = 3
}

variable "hpa_max_replicas" {
  type    = number
  default = 10
}
