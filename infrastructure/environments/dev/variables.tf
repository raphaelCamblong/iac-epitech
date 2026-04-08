variable "project_id" {
  type = string
}

variable "github_repo_url" {
  type        = string
  description = "Repository URL for ARC"
  default     = "https://github.com/raphaelCamblong/iac-epitech"
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
  type        = number
  default     = 1
  description = "GKE primary pool: nodes per zone in each configured zone. With the default single-zone pool (gke_node_locations), 1 is one worker VM."
}

variable "gke_node_locations" {
  type        = list(string)
  default     = null
  description = "Primary pool zones. null defaults to [region]-a only (one zone, lower cost). Set all regional zones for multi-AZ dev, e.g. [\"europe-west9-a\", \"europe-west9-b\", \"europe-west9-c\"]."
  nullable    = true
}

variable "node_machine_type" {
  type    = string
  default = "e2-standard-2"
}

variable "node_disk_size_gb" {
  type        = number
  description = "Boot disk per node; pd-standard (see GKE module) avoids SSD quota."
  default     = 20
}

variable "db_instance_name" {
  type    = string
  default = "task-manager-db"
}

variable "db_tier" {
  type    = string
  default = "db-f1-micro"
}

variable "db_availability_type" {
  type    = string
  default = "ZONAL"
}

variable "db_disk_size_gb" {
  type        = number
  description = "Cloud SQL allocation; 10 GB is the usual minimum for PostgreSQL."
  default     = 10
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
