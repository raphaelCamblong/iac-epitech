output "database_private_ip" {
  value = google_sql_database_instance.main.private_ip_address
}

output "database_name" {
  value = google_sql_database.app.name
}

output "database_user" {
  value = google_sql_user.app.name
}

output "database_password" {
  value = random_password.database.result
  sensitive = true
}

output "database_url" {
  value = "postgres://${google_sql_user.app.name}:${random_password.database.result}@${google_sql_database_instance.main.private_ip_address}:5432/${google_sql_database.app.name}?sslmode=disable"
  sensitive = true
}

output "connection_name" {
  value = google_sql_database_instance.main.connection_name
}
