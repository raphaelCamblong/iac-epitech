resource "random_password" "database" {
  length  = 24
  special = false
}

resource "google_sql_database_instance" "main" {
  name             = var.db_instance_name
  database_version = "POSTGRES_16"
  region           = var.region

  deletion_protection = false

  settings {
    tier                  = var.db_tier
    edition               = "ENTERPRISE"
    availability_type     = var.db_availability_type
    disk_size             = var.db_disk_size_gb
    disk_type             = var.db_disk_type
    disk_autoresize       = var.db_disk_autoresize
    disk_autoresize_limit = var.db_disk_autoresize_limit_gb

    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
    }

    ip_configuration {
      ipv4_enabled                                  = false
      private_network                               = var.network_id
      enable_private_path_for_google_cloud_services = true
    }
  }

  depends_on = [
    var.private_vpc_connection
  ]
}

resource "google_sql_database" "app" {
  name     = var.db_name
  instance = google_sql_database_instance.main.name
}

resource "google_sql_user" "app" {
  name     = var.db_user
  instance = google_sql_database_instance.main.name
  password = random_password.database.result
}
