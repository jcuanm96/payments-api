resource "google_cloudbuild_trigger" "api-build-trigger" {
  name = var.trigger_name
  description = "Build and deploy to Cloud Run service on push"

  github {
    owner = "VamaSingapore"
    name = "vama-api"
    push {
      branch = var.branch_regex
    }
  }

  substitutions = {
    _SERVICE_NAME = var.service_name
  }

  filename = "cloudbuild.yaml"
}