variable sql_instance {}

variable service_name {}

variable region {
  default = "us-central1"
}

variable zone {
  default = "us-central1-c"
}

variable port {
  default = "8080"
}

variable memory {
  default = "256Mi"
}

variable max_instances {
  default = "100"
}

variable request_timeout_seconds {
  default = "60"
}

terraform {
  experiments = [module_variable_optional_attrs]
}

variable env {
  type = set(
    object({
      key = string,
      value = optional(string),
      secret = optional(string),
    }),
  )

  default = []
  validation {
    error_message = "Environment variables must have one of `value` or `secret` defined."
    condition = alltrue([
      length([for e in var.env: e if (e.value == null && e.secret == null)]) < 1,
      length([for e in var.env: e if (e.value != null && e.secret != null)]) < 1,
    ])
  }
}