resource "google_service_account" "gke_nodes" {
  account_id   = "${var.name_prefix}-gke-nodes"
  display_name = "Task Manager GKE nodes"
}

locals {
  node_service_account_roles = toset([
    "roles/artifactregistry.reader",
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/stackdriver.resourceMetadata.writer",
  ])
}

resource "google_project_iam_member" "gke_nodes" {
  for_each = local.node_service_account_roles

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.gke_nodes.email}"
}

resource "google_container_cluster" "main" {
  name     = var.cluster_name
  location = var.cluster_location != null ? var.cluster_location : var.region

  network    = var.network_id
  subnetwork = var.subnetwork_id

  deletion_protection      = false
  remove_default_node_pool = true
  initial_node_count       = 1
  networking_mode          = "VPC_NATIVE"

  # Default pool exists only until remove_default_node_pool finishes. Without this block,
  # GKE uses API defaults (often pd-balanced) → SSD_TOTAL_GB quota even though the real pool uses pd-standard below.
  node_config {
    machine_type    = var.node_machine_type
    disk_size_gb    = var.node_disk_size_gb
    disk_type       = "pd-standard"
    service_account = google_service_account.gke_nodes.email
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]
  }

  release_channel {
    channel = var.gke_release_channel
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = var.pods_secondary_range_name
    services_secondary_range_name = var.services_secondary_range_name
  }

  lifecycle {
    ignore_changes = [
      node_config,
    ]
  }

  timeouts {
    create = "40m"
    update = "40m"
    delete = "40m"
  }
}

resource "google_container_node_pool" "primary" {
  name     = "${var.name_prefix}-pool"
  cluster  = google_container_cluster.main.name
  location = google_container_cluster.main.location

  node_locations     = var.node_locations
  initial_node_count = var.node_count

  autoscaling {
    min_node_count  = var.node_count
    max_node_count  = var.node_pool_max_count
    location_policy = "BALANCED"
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  lifecycle {
    ignore_changes = [
      initial_node_count,
    ]
  }

  node_config {
    machine_type    = var.node_machine_type
    disk_size_gb    = var.node_disk_size_gb
    disk_type       = "pd-standard"
    service_account = google_service_account.gke_nodes.email
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]

    labels = {
      app = "task-manager"
    }

    tags = [
      var.name_prefix,
      "gke-node",
    ]
  }
}
