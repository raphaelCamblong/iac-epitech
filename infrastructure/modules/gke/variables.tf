variable "project_id" {
  type = string
}

variable "region" {
  type = string
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
  type = number
}

variable "node_machine_type" {
  type = string
}

variable "node_disk_size_gb" {
  type = number
}
