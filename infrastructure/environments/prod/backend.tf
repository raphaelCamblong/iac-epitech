terraform {
  backend "gcs" {
    bucket = "terraform-state-task-manager-prod-492216" # must be create in gcp check the DEPLOYMENT.md
    prefix = "env/prod"
  }
}
