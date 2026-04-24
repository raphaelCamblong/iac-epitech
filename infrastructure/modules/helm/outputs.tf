output "ingress_nginx_status" {
  value = helm_release.ingress_nginx.status
}

output "task_manager_status" {
  value = try(helm_release.task_manager[0].status, "not deployed")
}
