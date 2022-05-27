# GCS Permissions
resource "google_storage_default_object_access_control" "public_rule" {
  bucket = google_storage_bucket.bucket.name
  role   = "READER"
  entity = "allUsers"
}


# GCS Buckets
resource "google_storage_bucket" "bucket" {
  name     = join("-", [var.bucket_name, var.gcloud_project])
  location = var.region
  uniform_bucket_level_access = var.uniform_bucket_level_access
}