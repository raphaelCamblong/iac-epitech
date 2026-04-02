variable "github_repo_url" {
  description = "The repository URL to configure ARC runners for"
  type        = string
}

variable "github_pat" {
  description = "GitHub Personal Access Token for ARC authentication"
  type        = string
  sensitive   = true
}
