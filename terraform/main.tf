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

resource "google_compute_network" "vpc_network" {
  name                    = "daily-dashboard-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "daily-dashboard-subnet"
  network       = google_compute_network.vpc_network.id
  ip_cidr_range = "10.0.0.0/24"
  region        = local.json_config["GCPRegion"]
}

resource "google_compute_instance" "default" {
  name         = "daily-dashboard"
  machine_type = "f1-micro"
  zone         = local.json_config["GCPZone"]
  tags         = ["ssh"]

  boot_disk {
    initialize_params {
      image = "ubuntu-2204-lts"
    }
  }

  metadata_startup_script = "sudo apt update"
  network_interface {
    subnetwork = google_compute_subnetwork.default.id
    access_config {}
  }
}

resource "google_compute_firewall" "ssh" {
  name = "allow-ssh"
  allow {
    ports    = ["22"]
    protocol = "tcp"
  }
  direction     = "INGRESS"
  network       = google_compute_network.vpc_network.id
  priority      = 1000
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["ssh"]
}

resource "google_compute_firewall" "daily-dashboard" {
  name    = "daily-dashboard-firewall"
  network = google_compute_network.vpc_network.id

  allow {
    protocol = "tcp"
    ports    = ["8080"]
  }
  source_ranges = ["0.0.0.0/0"]
}

output "Web-server-URL" {
  value = join("", ["http://", google_compute_instance.default.network_interface.0.access_config.0.nat_ip, ":8080"])
}
