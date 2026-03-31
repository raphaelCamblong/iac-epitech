terraform {
  backend "gcs" {
    bucket = "terraform-state-task-manager-prod" # You must create this bucket
    prefix = "env/prod"
  }
}
