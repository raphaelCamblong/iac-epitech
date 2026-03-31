variable "platform_state_path" {
  description = "Optional local path to the platform Terraform state file."
  type        = string
  default     = "platform/terraform.tfstate"
}

variable "project_id" {
  description = "GCP project id. Leave empty when using platform_state_path."
  type        = string
  default     = ""
}

variable "region" {
  description = "GCP region. Leave empty when using platform_state_path."
  type        = string
  default     = ""
}

variable "cluster_name" {
  description = "GKE cluster name. Leave empty when using platform_state_path."
  type        = string
  default     = ""
}

variable "cluster_location" {
  description = "GKE cluster location. Leave empty when using platform_state_path."
  type        = string
  default     = ""
}

variable "namespace" {
  description = "Kubernetes namespace for the deployment"
  type        = string
  default     = "default"
}

variable "release_name" {
  description = "Helm release name"
  type        = string
  default     = "task-manager"
}

variable "database_url" {
  description = "PostgreSQL connection URL. Leave empty when using platform_state_path."
  type        = string
  default     = ""
  sensitive   = true
}

variable "jwt_secret" {
  description = "JWT signing secret"
  type        = string
  sensitive   = true
}

variable "image_repository" {
  description = "Full container image repository. If empty, it is built from artifact_registry_repository + image_name."
  type        = string
  default     = ""
}

variable "artifact_registry_repository" {
  description = "Artifact Registry base path. Leave empty when using platform_state_path."
  type        = string
  default     = ""
}

variable "image_name" {
  description = "Docker image name appended to artifact_registry_repository when image_repository is empty."
  type        = string
  default     = "task-manager"
}

variable "image_tag" {
  description = "Container image tag"
  type        = string
  default     = "latest"
}

variable "ingress_host" {
  description = "Ingress hostname for the API"
  type        = string
  default     = "task-manager.example.com"
}

variable "ingress_static_ip" {
  description = "Optional reserved static IP used by the ingress controller service."
  type        = string
  default     = ""
}

variable "replica_count" {
  description = "Number of replicas"
  type        = number
  default     = 2
}

variable "hpa_min_replicas" {
  description = "HPA minimum replicas"
  type        = number
  default     = 2
}

variable "hpa_max_replicas" {
  description = "HPA maximum replicas"
  type        = number
  default     = 10
}

variable "ingress_controller_namespace" {
  description = "Namespace used by ingress-nginx."
  type        = string
  default     = "ingress-nginx"
}

variable "ingress_nginx_chart_version" {
  description = "ingress-nginx chart version."
  type        = string
  default     = "4.11.2"
}
