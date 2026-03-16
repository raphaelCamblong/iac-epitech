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
