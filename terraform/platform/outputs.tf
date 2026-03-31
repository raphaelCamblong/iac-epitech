output "project_id" {
  description = "GCP project id used by the platform."
  value       = var.project_id
}

output "region" {
  description = "GCP region used by the platform."
  value       = var.region
}

output "cluster_name" {
  description = "GKE cluster name."
  value       = google_container_cluster.main.name
}

output "cluster_location" {
  description = "GKE cluster location."
  value       = google_container_cluster.main.location
}

output "artifact_registry_repository" {
  description = "Base Artifact Registry path used to push images."
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.images.repository_id}"
}

output "ingress_static_ip" {
  description = "Reserved regional IP for the ingress controller."
  value       = google_compute_address.ingress.address
}

output "database_private_ip" {
  description = "Private IP of the Cloud SQL instance."
  value       = google_sql_database_instance.main.private_ip_address
}

output "database_name" {
  description = "Application database name."
  value       = google_sql_database.app.name
}

output "database_user" {
  description = "Application database user."
  value       = google_sql_user.app.name
}

output "database_password" {
  description = "Application database password."
  value       = random_password.database.result
  sensitive   = true
}

output "database_url" {
  description = "Connection string consumed by the application."
  value       = "postgres://${google_sql_user.app.name}:${random_password.database.result}@${google_sql_database_instance.main.private_ip_address}:5432/${google_sql_database.app.name}?sslmode=disable"
  sensitive   = true
}

output "state_bucket_name" {
  description = "Bucket created to host Terraform state later if desired."
  value       = google_storage_bucket.terraform_state.name
}
