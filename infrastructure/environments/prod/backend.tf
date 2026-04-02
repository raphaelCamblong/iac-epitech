terraform {
  backend "gcs" {
    bucket = "terraform-state-task-manager-prod" # must be create in gcp check the DEPLOYMENT.md
    prefix = "env/prod"
  }
}
