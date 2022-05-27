resource "google_cloud_tasks_queue" "queue" {
  name = var.instance_name
  location = var.region
  project = var.gcloud_project

  rate_limits {
    max_concurrent_dispatches = 1
    max_dispatches_per_second = 2
  }

  retry_config {
    max_attempts = -1
    max_backoff = "10800s" // 3 hours
    min_backoff = "2s"
    max_doublings = 10
  }
}