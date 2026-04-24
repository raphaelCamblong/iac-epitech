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

variable "gke_cluster_location" {
  type        = string
  default     = null
  description = "Cluster location. null → regional cluster in var.region (full HA control plane, multi-AZ default pool). Set a zone (e.g. europe-west9-a) to fall back to a zonal cluster if regional capacity/quota is exhausted."
  nullable    = true
}

variable "node_count" {
  type        = number
  description = "GKE primary pool: nodes per zone. Prod uses all regional zones, so total VMs ≈ node_count × zone count (often three)."
  default     = 1
}

variable "gke_node_pool_max_count" {
  type        = number
  description = "Max nodes per zone (autoscaling). Keeps total disks bounded: ≈ max × zones × node_disk_size_gb."
  default     = 2
}

variable "node_machine_type" {
  type        = string
  description = "Same light default as dev."
  default     = "e2-standard-2"
}

variable "node_disk_size_gb" {
  type        = number
  description = "Boot disk per node (pd-standard in GKE module = standard PD quota, not SSD). Prod is multi-zone: total node disks ≈ size × node_count × zones."
  default     = 20
}

variable "db_instance_name" {
  type    = string
  default = "task-manager-db"
}

variable "db_tier" {
  type        = string
  description = "Right-sized for a small Postgres workload; override for more CPU/RAM."
  default     = "db-f1-micro"
}

variable "db_availability_type" {
  type        = string
  description = "REGIONAL duplicates the instance (higher cost). ZONAL is enough for many small apps."
  default     = "ZONAL"
}

variable "db_disk_size_gb" {
  type        = number
  description = "Cloud SQL initial allocation (GB). Increase via Terraform if needed; shrinking often requires instance replacement."
  default     = 10
}

variable "db_disk_type" {
  type        = string
  description = "PD_HDD recommended for cost and SSD quota; PD_SSD counts against SSD_TOTAL_GB."
  default     = "PD_HDD"
}

variable "db_disk_autoresize" {
  type        = bool
  description = "Let Cloud SQL grow disk automatically up to db_disk_autoresize_limit_gb."
  default     = true
}

variable "db_disk_autoresize_limit_gb" {
  type        = number
  description = "Cap autoresize so a DB cannot grow into hundreds of GB by accident."
  default     = 30
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
  default = 2
}

variable "hpa_min_replicas" {
  type    = number
  default = 2
}

variable "hpa_max_replicas" {
  type    = number
  default = 10
}

variable "deploy_app" {
  type        = bool
  default     = true
  description = "Whether to install the task-manager Helm release. Set to false (via tfvars or -var) on the very first apply, before CI has pushed an image to Artifact Registry. Once an image exists, flip to true (the default) for normal CD runs."
}
