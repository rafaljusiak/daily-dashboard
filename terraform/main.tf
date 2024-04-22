terraform {
  required_providers {
    external = {
      source  = "hashicorp/external"
      version = "2.3.3"
    }
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  project = local.json_config["GCPProjectId"]
}

resource "google_compute_network" "vpc_network" {
  name = "daily-dashboard-network"
}

locals {
  json_config = jsondecode(file("${path.module}/../config.json"))
}
