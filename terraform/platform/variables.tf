variable "project_id" {
  description = "Google Cloud project id."
  type        = string
}

variable "region" {
  description = "Primary GCP region."
  type        = string
  default     = "europe-west1"
}

variable "zone" {
  description = "Primary GCP zone."
  type        = string
  default     = "europe-west1-b"
}

variable "name_prefix" {
  description = "Short prefix used for GCP resource names."
  type        = string
  default     = "task-manager"
}

variable "cluster_name" {
  description = "GKE cluster name."
  type        = string
  default     = "task-manager-gke"
}

variable "gke_release_channel" {
  description = "GKE release channel."
  type        = string
  default     = "REGULAR"
}

variable "node_count" {
  description = "Number of nodes in the primary node pool."
  type        = number
  default     = 2
}

variable "node_machine_type" {
  description = "Machine type used by the GKE nodes."
  type        = string
  default     = "e2-standard-2"
}

variable "node_disk_size_gb" {
  description = "Disk size for each GKE node."
  type        = number
  default     = 50
}

variable "subnet_cidr" {
  description = "Primary subnet CIDR used by the cluster nodes."
  type        = string
  default     = "10.10.0.0/20"
}

variable "pods_cidr" {
  description = "Secondary CIDR used by Kubernetes pods."
  type        = string
  default     = "10.20.0.0/16"
}

variable "services_cidr" {
  description = "Secondary CIDR used by Kubernetes services."
  type        = string
  default     = "10.30.0.0/20"
}

variable "pods_secondary_range_name" {
  description = "Secondary range name for pods."
  type        = string
  default     = "pods"
}

variable "services_secondary_range_name" {
  description = "Secondary range name for services."
  type        = string
  default     = "services"
}

variable "db_instance_name" {
  description = "Cloud SQL instance name."
  type        = string
  default     = "task-manager-db"
}

variable "db_name" {
  description = "Application database name."
  type        = string
  default     = "taskmanager"
}

variable "db_user" {
  description = "Application database user."
  type        = string
  default     = "taskmanager"
}

variable "db_tier" {
  description = "Cloud SQL machine tier."
  type        = string
  default     = "db-custom-1-3840"
}

variable "db_disk_size_gb" {
  description = "Cloud SQL disk size."
  type        = number
  default     = 20
}

variable "db_availability_type" {
  description = "Cloud SQL availability type."
  type        = string
  default     = "ZONAL"
}

variable "artifact_registry_name" {
  description = "Artifact Registry repository name."
  type        = string
  default     = "task-manager"
}

variable "state_bucket_name" {
  description = "Optional custom bucket name for Terraform state artifacts."
  type        = string
  default     = ""
}
