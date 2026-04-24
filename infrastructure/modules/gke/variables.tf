variable "project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "cluster_location" {
  type        = string
  default     = null
  description = "Cluster location. Zone (e.g. europe-west9-a) → zonal cluster (1 control plane, default pool = initial_node_count VMs). null → regional cluster in var.region (3 control planes, default pool = initial_node_count × 3 VMs, needs capacity in every zone)."
  nullable    = true
}

variable "name_prefix" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "network_id" {
  type = string
}

variable "subnetwork_id" {
  type = string
}

variable "gke_release_channel" {
  type = string
}

variable "pods_secondary_range_name" {
  type = string
}

variable "services_secondary_range_name" {
  type = string
}

variable "node_count" {
  type        = number
  description = "Nodes per zone for each zone in node_locations. If node_locations is unset, the pool uses every zone in the regional cluster (typically three), so 1 can mean three VMs total."
}

variable "node_locations" {
  type        = list(string)
  default     = null
  description = "Zones for the primary node pool (e.g. [\"europe-west9-a\"]). Null omits the field so GKE uses all zones in the cluster region."
  nullable    = true
}

variable "node_machine_type" {
  type = string
}

variable "node_disk_size_gb" {
  type = number
}

variable "node_pool_max_count" {
  type        = number
  default     = 5
  description = "Autoscaling ceiling per zone. In a regional cluster, worst-case node count ≈ this × number of zones."
}
