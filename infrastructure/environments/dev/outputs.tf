output "ingress_url" {
  value = "http://${module.network.ingress_static_ip}.nip.io"
}
