
resource "google_storage_bucket" "bucket" {
  name     = var.bucketName
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = var.sourceCodeZipName
  bucket = google_storage_bucket.bucket.name
  source = var.sourceCodePath
}

# Cloud function does NOT automatically detect changes to
# the storage bucket object & redeploy. If you change the underlying object,
# delete the existing Cloud Function before running `terraform plan`.
# https://github.com/hashicorp/terraform-provider-google/issues/1938
resource "google_cloudfunctions_function" "function" {
  name        = var.functionName
  description = var.functionDescription
  runtime     = "python39"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  entry_point = var.entryPoint

  environment_variables = {
    VAMA_API_BASE_URL = var.baseURL
    TELEGRAM_HEALTH_CHECKER_API_TOKEN = var.telegramToken
  }
}

# IAM entry for all users to invoke the function
resource "google_cloudfunctions_function_iam_member" "invoker" {
  project        = google_cloudfunctions_function.function.project
  region         = google_cloudfunctions_function.function.region
  cloud_function = google_cloudfunctions_function.function.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}