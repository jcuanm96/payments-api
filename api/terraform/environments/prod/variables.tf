# GCP Env
variable "gcloud_project" {
  default = "vama-prod"
}

variable "api_base_url" {
  default = "https://api-prod-x7dqcpkzxa-uc.a.run.app"
}

variable "service_account_email" {
  default = "vama-prod@vama-prod.iam.gserviceaccount.com"
}

variable "credentials_file" {
  default = "creds-prod.json"
}

variable "region" {
  default = "us-central1"
}

variable "zone" {
  default = "us-central1-c"
}

# Cloud Run
variable "service_name" {
  default = "api-prod"
}

variable "logger_type" {
  default = "GCP"
}