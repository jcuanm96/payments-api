resource "google_sql_database_instance" "instance" {
  name             = var.db_name
  region           = var.region
  database_version = "POSTGRES_13"
  settings {
    tier = "db-g1-small"
    maintenance_window {
      day  = 2
      hour = 12
    }
  }
}