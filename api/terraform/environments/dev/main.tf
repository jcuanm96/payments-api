provider "google" {
  credentials = file(var.credentials_file)

  project = var.gcloud_project
  region  = var.region
  zone    = var.zone
}

module "db" {
  source  = "../../modules/DbInstance"
  db_name = "vama-dev-postgres"
}

module "cloud_build" {
  source       = "../../modules/CloudBuild"
  trigger_name = "dev-build-trigger"
  branch_regex = "^dev$"
  service_name = var.service_name
}

module "cloud_run" {
  source       = "../../modules/CloudRun"
  sql_instance = module.db.instance_connection_name
  service_name = var.service_name
  memory       = "1Gi"

  // Environment variables
  env = [
    { key = "POSTGRE_CONN_STR", secret = "postgre_conn_str" },

    { key = "TWILIO_SID", secret = "twilio_sid" },
    { key = "TWILIO_TOKEN", secret = "twilio_token" },
    { key = "TWILIO_VERIFY", secret = "twilio_verify" },

    { key = "STRIPE_KEY", secret = "stripe_api_key" },
    { key = "STRIPE_EVENT_SECRET", secret = "stripe_event_secret" },

    { key = "SENDBIRD_MASTER_API_KEY", secret = "sendbird_master_api_key" },
    { key = "SENDBIRD_APPLICATION_ID", secret = "sendbird_application_id" },

    { key = "ACCESS_TOKEN_KEY", secret = "access_token_key" },
    { key = "REFRESH_TOKEN_KEY", secret = "refresh_token_key" },
    { key = "GOAT_INVITE_CODE_SECRET", secret = "goat_invite_code_secret" },
    { key = "PAYOUT_SECRET", secret = "payout_secret" },
    { key = "DASHBOARD_SECRET", secret = "dashboard_secret" },
    { key = "GCLOUD_PROJECT", value = var.gcloud_project },
    { key = "SERVICE_ACCOUNT_EMAIL", value = var.service_account_email },
    { key = "THEME_SECRET", secret = "theme_secret" },
    { key = "API_BASE_URL", value = var.api_base_url },
    { key = "REDIRECT_BASE_URL", value = "localhost:5001" },
    { key = "TELEGRAM_HEALTH_CHECKER_API_TOKEN", value = "5084315086:AAF26WWMJ41C6oc-F0pDTG7jMV0GbMc8atg" },
    { key = "LOGGER_TYPE", value = var.logger_type },
    { key = "REDIS_ENDPOINT", value = "redis-15524.c18564.us-central1-1.gcp.cloud.rlrcp.com" },
    { key = "REDIS_PORT", value = "15524" },
    { key = "REDIS_PASSWORD", secret = "vama_api_redis_password" }
  ]
}

module "cloud_function" {
  source              = "../../modules/CloudFunctions"
  sourceCodePath      = "health.zip"
  sourceCodeZipName   = "health.zip"
  bucketName          = "vama_health_checker_dev_bucket"
  functionName        = "vama_health_checker_dev"
  functionDescription = "Vama Health Checker Function for dev"
  baseURL             = module.cloud_run.url
  telegramToken       = "5084315086:AAF26WWMJ41C6oc-F0pDTG7jMV0GbMc8atg"
  entryPoint          = "check_health"
}

module "cloud_scheduler" {
  source               = "../../modules/CloudScheduler"
  url                  = module.cloud_function.url
  schedule             = "*/5 * * * *" // Every 5th minute
  schedulerName        = "vama_health_scheduler_dev"
  schedulerDescription = "Vama Health Checker Scheduler for dev"
  httpMethod           = "POST"
  serviceAccountEmail  = var.service_account_email
  audience             = module.cloud_run.url
}

module "pending_balance_notifications_cloud_scheduler" {
  source               = "../../modules/CloudScheduler"
  url                  = format("%s/%s", module.cloud_run.url, "cloudscheduler/v1/wallet/balance/notify")
  schedule             = "0 21 * * 1" // At 21:00 UTC on every Monday
  schedulerName        = "pending-balance-notifications"
  schedulerDescription = "Vama Pending Balance Push Notifications"
  httpMethod           = "POST"
  serviceAccountEmail  = var.service_account_email
  audience             = module.cloud_run.url
}

module "profile_avatars_gcs_bucket" {
  source         = "../../modules/GCS"
  gcloud_project = var.gcloud_project
  bucket_name    = "profile-avatars"
}

module "feed_posts_gcs_bucket" {
  source         = "../../modules/GCS"
  gcloud_project = var.gcloud_project
  bucket_name    = "feed-posts"
}

# Doesn't use GCS module to make this a fully private bucket 
resource "google_storage_bucket" "bucket" {
  name                        = join("-", ["chat-media", var.gcloud_project])
  location                    = var.region
}
resource "google_storage_default_object_access_control" "public_rule" {
  bucket = google_storage_bucket.bucket.name
  role   = "READER"
  entity = "allUsers"
}

module "stripe_unsubscribe_cloud_task_queue" {
  source         = "../../modules/CloudTask"
  gcloud_project = var.gcloud_project
  instance_name  = "stripe-paid-group-unsubscribe"
}

module "remove_from_paid_group_cloud_task_queue" {
  source         = "../../modules/CloudTask"
  gcloud_project = var.gcloud_project
  instance_name  = "remove-from-paid-group"
}

module "add_user_contacts_task_queue" {
  source         = "../../modules/CloudTask"
  gcloud_project = var.gcloud_project
  instance_name  = "add-user-contacts"
}

# Redis resources
data "rediscloud_payment_method" "card" {
  card_type = "Visa"
}

resource "random_password" "api_redis_password" {
 length = 20
 upper = true
 lower = true
 number = true
 special = false
}
 
resource "rediscloud_subscription" "vama_redis_subscription" {
  name = join("-", ["api", var.gcloud_project])
  memory_storage = "ram"
  payment_method_id = data.rediscloud_payment_method.card.id
  persistent_storage_encryption = false

  cloud_provider {
    provider = "GCP"
    region {
      region = "us-central1"
      networking_deployment_cidr = "192.168.0.0/24"
      preferred_availability_zones = ["us-central1-a"]
    }
  }

  database {
    name = join("-", ["api", var.gcloud_project])
    protocol = "redis"
    memory_limit_in_gb = 1
    replication = true
    data_persistence = "none"
    throughput_measurement_by = "number-of-shards"
    throughput_measurement_value = 2
    password = random_password.api_redis_password.result
  }
}
