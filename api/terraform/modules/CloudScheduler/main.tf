resource "google_cloud_scheduler_job" "job" {
  name             = var.schedulerName
  description      = var.schedulerDescription

  schedule         = var.schedule
  attempt_deadline = "15s"

  retry_config {
    retry_count = 0
  }

  http_target {
    http_method = var.httpMethod
    uri         = var.url
    body        = base64encode("{}")

    oidc_token {
      service_account_email = var.serviceAccountEmail
      audience = var.audience
    }
  }
}