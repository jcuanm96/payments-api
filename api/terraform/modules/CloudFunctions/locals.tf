locals {
  env = toset([
    for e in var.env: {
      key = e.key
      value = e.value
      secret = {
        name = e.secret
        version = "latest"
      }
    }
  ])
}