terraform {
  backend "gcs" {
    bucket = "terraform-state-task-manager-dev" # You must create this bucket
    prefix = "env/dev"
  }
}
