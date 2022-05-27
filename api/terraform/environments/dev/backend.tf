terraform {
  backend "gcs" {
    bucket      = "vama-dev-tfstate"
    prefix      = "env/dev"
    credentials = "creds-dev.json"
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
