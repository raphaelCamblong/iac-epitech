variable "kube_config_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
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
  description = "PostgreSQL connection URL (from CloudSQL, RDS, or external)"
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "JWT signing secret"
  type        = string
  sensitive   = true
}

variable "image_repository" {
  description = "Container image repository"
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
