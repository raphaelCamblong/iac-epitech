output "network_id" {
  value = google_compute_network.main.id
}

output "subnetwork_id" {
  value = google_compute_subnetwork.gke.id
}

output "ingress_static_ip" {
  value = google_compute_address.ingress.address
}

output "private_vpc_connection" {
  value = google_service_networking_connection.private_vpc_connection.id
}
