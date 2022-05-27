# GCP Env
variable "gcloud_project" {
  default = "vama-dev"
}

variable "api_base_url" {
  default = "https://api-dev-onxum5hzma-uc.a.run.app"
}

variable "service_account_email" {
  default = "vama-dev@vama-dev.iam.gserviceaccount.com"
}

variable "credentials_file" {
  default = "creds-dev.json"
}

variable "region" {
  default = "us-central1"
}

variable "zone" {
  default = "us-central1-c"
}

# Cloud Run
variable "service_name" {
  default = "api-dev"
}

variable "logger_type" {
  default = "GCP"
}