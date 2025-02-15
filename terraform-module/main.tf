terraform {
  required_providers {
    carstore = {
      source = "local/carstore"
      version = "1.0.0"
    }
  }
}

provider "carstore" {
  base_url = "http://localhost:5000"

}

resource "carstore_car" "my_car" {
  model = "Tesla"
  year  = 2024
}