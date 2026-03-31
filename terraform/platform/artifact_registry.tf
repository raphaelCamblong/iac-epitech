resource "google_artifact_registry_repository" "images" {
  location      = var.region
  repository_id = var.artifact_registry_name
  description   = "Docker images for Task Manager"
  format        = "DOCKER"

  depends_on = [google_project_service.required["artifactregistry.googleapis.com"]]
}
