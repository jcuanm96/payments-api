# GCP Env
variable "gcloud_project" {
  default = "vama-staging"
}

variable "api_base_url" {
  default = "https://api-staging-ypidedhdqa-uc.a.run.app"
}

variable "service_account_email" {
  default = "vama-staging@vama-staging.iam.gserviceaccount.com"
}

variable "credentials_file" {
  default = "creds-staging.json"
}

variable "region" {
  default = "us-central1"
}

variable "zone" {
  default = "us-central1-c"
}

# Cloud Run
variable "service_name" {
  default = "api-staging"
}

variable "logger_type" {
  default = "GCP"
}