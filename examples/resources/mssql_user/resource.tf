resource "mssql_login" "example" {
  name = "example"
  password = "example"
}

resource "mssql_user" "example" {
  name = mssql_login.example.name
  database = "example-db"
}
