terraform {
  required_providers {
    mssql = {
      source = "terraform.embracecloud.nl/embracecloud/mssql"
    }
  }
}


provider "mssql" {
  connection_string = "localhost:44444333"
}

resource "mssql_login" "name" {
  name = "test"
  password = "test"
}