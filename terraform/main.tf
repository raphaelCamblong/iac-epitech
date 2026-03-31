terraform {
  required_version = ">= 1.6.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

data "terraform_remote_state" "platform" {
  count = var.platform_state_path == "" ? 0 : 1

  backend = "local"

  config = {
    path = var.platform_state_path
  }
}

locals {
  platform_outputs = var.platform_state_path != "" ? data.terraform_remote_state.platform[0].outputs : {}

  project_id = var.project_id != "" ? var.project_id : try(local.platform_outputs.project_id, "")
  region     = var.region != "" ? var.region : try(local.platform_outputs.region, "")

  cluster_name     = var.cluster_name != "" ? var.cluster_name : try(local.platform_outputs.cluster_name, "")
  cluster_location = var.cluster_location != "" ? var.cluster_location : try(local.platform_outputs.cluster_location, "")

  database_url = var.database_url != "" ? var.database_url : try(local.platform_outputs.database_url, "")

  artifact_registry_repository = var.artifact_registry_repository != "" ? var.artifact_registry_repository : try(local.platform_outputs.artifact_registry_repository, "")
  image_repository             = var.image_repository != "" ? var.image_repository : (local.artifact_registry_repository != "" ? "${local.artifact_registry_repository}/${var.image_name}" : "")

  ingress_static_ip = var.ingress_static_ip != "" ? var.ingress_static_ip : try(local.platform_outputs.ingress_static_ip, "")
}

provider "google" {
  project = local.project_id
  region  = local.region
}

data "google_client_config" "current" {}

data "google_container_cluster" "target" {
  name     = local.cluster_name
  location = local.cluster_location
  project  = local.project_id
}

provider "kubernetes" {
  host                   = "https://${data.google_container_cluster.target.endpoint}"
  token                  = data.google_client_config.current.access_token
  cluster_ca_certificate = base64decode(data.google_container_cluster.target.master_auth[0].cluster_ca_certificate)
}

provider "helm" {
  kubernetes {
    host                   = "https://${data.google_container_cluster.target.endpoint}"
    token                  = data.google_client_config.current.access_token
    cluster_ca_certificate = base64decode(data.google_container_cluster.target.master_auth[0].cluster_ca_certificate)
  }
}
