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

variable sourceCodePath {}

variable entryPoint {}

variable functionName {}

variable functionDescription {}

variable sourceCodeZipName {}

variable bucketName {}

variable baseURL {}

variable telegramToken {}