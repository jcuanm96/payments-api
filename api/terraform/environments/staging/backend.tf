terraform {
  backend "gcs" {
    bucket      = "vama-staging-tfstate"
    prefix      = "env/staging"
    credentials = "creds-staging.json"
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.3.0"
    }
    rediscloud = {
      source = "RedisLabs/rediscloud"
      version = "0.2.9"
    }
  }
}