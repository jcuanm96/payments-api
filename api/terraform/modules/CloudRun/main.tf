
resource "google_cloud_run_service" "default" {
  provider = google
  name     = var.service_name
  location = var.region
  autogenerate_revision_name = true

  template {
    spec {
      timeout_seconds = var.request_timeout_seconds
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
        resources {
          limits = {
            memory = var.memory
            cpu = "1000m"
          }  
        }
        ports {
          name = "http1"
          container_port = var.port
        }

         # Populate straight environment variables.
        dynamic env {
          for_each = [for e in local.env: e if e.value != null]

          content {
            name = env.value.key
            value = env.value.value
          }
        }

        # Populate environment variables from secrets.
        dynamic env {
          for_each = [for e in local.env: e if e.secret.name != null]

          content {
            name = env.value.key
            value_from {
              secret_key_ref {
                name = env.value.secret.name
                key = env.value.secret.version
              }
            }
          }
        }
      }
    }
  
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = var.max_instances
        "autoscaling.knative.dev/minScale"      = "1"
        "run.googleapis.com/cloudsql-instances" = var.sql_instance
        "run.googleapis.com/client-name"        = "gcloud"
      }
    }
  }

   lifecycle {
    ignore_changes = [
      template[0].metadata[0].annotations["client.knative.dev/user-image"],
      template[0].metadata[0].annotations["run.googleapis.com/client-name"],
      template[0].metadata[0].annotations["run.googleapis.com/client-version"],
      template[0].metadata[0].annotations["run.googleapis.com/sandbox"],
      template[0].metadata[0].labels["commit-sha"],
      template[0].metadata[0].labels["gcb-build-id"],
      template[0].metadata[0].labels["gcb-trigger-id"],
      template[0].metadata[0].labels["managed-by"],
      template[0].spec[0].containers[0].image,
      template[0].spec[0].service_account_name,
      metadata[0].annotations["serving.knative.dev/creator"],
      metadata[0].annotations["serving.knative.dev/lastModifier"],
      metadata[0].annotations["run.googleapis.com/ingress-status"],
      metadata[0].labels["cloud.googleapis.com/location"],
    ]
  }

}

# This IAM policy allows unauthenticated HTTP requests to the 
# Cloud Run service.
data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = google_cloud_run_service.default.location
  project     = google_cloud_run_service.default.project
  service     = google_cloud_run_service.default.name

  policy_data = data.google_iam_policy.noauth.policy_data
} 