terraform {
  backend "gcs" {
    bucket = "terraform-state-task-manager-dev" # must be create in gcp check the DEPLOYMENT.md
    prefix = "env/dev"
  }
}
