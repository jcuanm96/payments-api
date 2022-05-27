resource "google_cloud_run_domain_mapping" "default" {
  location = var.region
  name     = var.domain

  metadata {
    namespace = var.namespace
  }

  spec {
    route_name = var.service_name
  }
}