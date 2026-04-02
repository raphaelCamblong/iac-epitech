output "artifact_registry_repository_url" {
  description = "Base Artifact Registry path used to push images."
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.images.repository_id}"
}
