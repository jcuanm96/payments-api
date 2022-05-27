variable "bucket_name" {}

variable "gcloud_project" {}

variable "region" {
  default = "us-central1"
}

variable "uniform_bucket_level_access" {
  default = false
}