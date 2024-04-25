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

locals {
  json_config = jsondecode(file("${path.module}/../config.json"))
}

provider "google" {
  project = local.json_config["GCPProjectId"]
}

resource "google_artifact_registry_repository" "dailydashboardrepository" {
  location      = local.json_config["GCPRegion"]
  repository_id = "daily-dashboard-repository"
  format        = "DOCKER"
}

resource "google_cloud_run_service" "dailydashboard" {
  name     = "dailydashboard"
  location = local.json_config["GCPRegion"]
  template {
    spec {
      containers {
        image = "${local.json_config["GCPRegion"]}-docker.pkg.dev/${local.json_config["GCPProjectId"]}/daily-dashboard-repository/daily-dashboard"
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.dailydashboard.location
  project  = google_cloud_run_service.dailydashboard.project
  service  = google_cloud_run_service.dailydashboard.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

output "url" {
  value = google_cloud_run_service.dailydashboard.status[0].url
}
