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
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Network Module
module "network" {
  source = "../../modules/network"

  name_prefix                   = "${var.name_prefix}-dev"
  region                        = var.region
  subnet_cidr                   = var.subnet_cidr
  pods_secondary_range_name     = var.pods_secondary_range_name
  pods_cidr                     = var.pods_cidr
  services_secondary_range_name = var.services_secondary_range_name
  services_cidr                 = var.services_cidr
}

# GKE Module
module "gke" {
  source = "../../modules/gke"

  project_id                    = var.project_id
  region                        = var.region
  name_prefix                   = "${var.name_prefix}-dev"
  cluster_name                  = "${var.name_prefix}-gke-dev"
  network_id                    = module.network.network_id
  subnetwork_id                 = module.network.subnetwork_id
  gke_release_channel           = var.gke_release_channel
  pods_secondary_range_name     = var.pods_secondary_range_name
  services_secondary_range_name = var.services_secondary_range_name
  node_count                    = var.node_count
  node_locations                = var.gke_node_locations != null ? var.gke_node_locations : ["${var.region}-a"]
  node_machine_type             = var.node_machine_type
  node_disk_size_gb             = var.node_disk_size_gb
}

# Database Module
module "database" {
  source = "../../modules/database"

  db_instance_name                      = "${var.db_instance_name}-dev"
  region                                = var.region
  db_tier                               = var.db_tier
  db_availability_type                  = var.db_availability_type
  db_disk_size_gb                       = var.db_disk_size_gb
  network_id                            = module.network.network_id
  private_service_networking_connection = module.network.private_service_networking_connection
  db_name                               = var.db_name
  db_user                               = var.db_user

  depends_on = [module.network]
}

# Artifact Registry Module
module "artifact_registry" {
  source = "../../modules/artifact_registry"

  region                 = var.region
  project_id             = var.project_id
  artifact_registry_name = "${var.artifact_registry_name}-dev"
}

data "google_client_config" "current" {}

provider "kubernetes" {
  host                   = "https://${module.gke.endpoint}"
  token                  = data.google_client_config.current.access_token
  cluster_ca_certificate = base64decode(module.gke.cluster_ca_certificate)
}

provider "helm" {
  kubernetes {
    host                   = "https://${module.gke.endpoint}"
    token                  = data.google_client_config.current.access_token
    cluster_ca_certificate = base64decode(module.gke.cluster_ca_certificate)
  }
}

# Helm (Application) Module
module "helm" {
  source = "../../modules/helm"

  ingress_static_ip = module.network.ingress_static_ip
  database_url      = module.database.database_url
  jwt_secret        = var.jwt_secret
  image_repository  = "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_registry_name}-dev/task-manager"
  image_tag         = var.image_tag
  ingress_host      = "${module.network.ingress_static_ip}.nip.io"
  replica_count     = var.replica_count
  hpa_min_replicas  = var.hpa_min_replicas
  hpa_max_replicas  = var.hpa_max_replicas
}

# ARC Module
module "arc" {
  source = "../../modules/arc"

  github_repo_url = var.github_repo_url
  github_pat      = var.github_pat

  depends_on = [module.gke]
}
