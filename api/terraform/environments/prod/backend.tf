terraform {
  backend "gcs" {
    bucket      = "vama-prod-tfstate"
    prefix      = "env/prod"
    credentials = "creds-prod.json"
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