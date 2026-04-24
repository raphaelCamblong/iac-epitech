variable "ingress_controller_namespace" {
  type    = string
  default = "ingress-nginx"
}

variable "ingress_nginx_chart_version" {
  type    = string
  default = "4.10.0"
}

variable "ingress_static_ip" {
  type = string
}

variable "namespace" {
  type    = string
  default = "default"
}

variable "release_name" {
  type    = string
  default = "task-manager"
}

variable "database_url" {
  type = string
}

variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "image_repository" {
  type = string
}

variable "image_tag" {
  type = string
}

variable "replica_count" {
  type    = number
  default = 2
}

variable "ingress_host" {
  type = string
}

variable "hpa_min_replicas" {
  type    = number
  default = 2
}

variable "hpa_max_replicas" {
  type    = number
  default = 5
}

variable "deploy_app" {
  type        = bool
  default     = true
  description = "Install the task-manager Helm release and its secret. Set to false on the first-ever apply (bootstrap), before CI has pushed an image to Artifact Registry, to avoid helm_release hanging on ImagePullBackOff. Flip back to true for normal CD runs."
}
