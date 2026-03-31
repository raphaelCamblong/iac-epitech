output "release_name" {
  description = "Helm release name"
  value       = helm_release.task_manager.name
}

output "namespace" {
  description = "Kubernetes namespace"
  value       = helm_release.task_manager.namespace
}

output "ingress_host" {
  description = "Ingress host for the API"
  value       = var.ingress_host
}

output "image_repository" {
  description = "Container repository deployed by Helm."
  value       = local.image_repository
}

output "ingress_static_ip" {
  description = "Static IP exposed by ingress-nginx when provided."
  value       = local.ingress_static_ip
}
