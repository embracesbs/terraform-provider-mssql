resource "mssql_login" "example" {
  name = "example"
  password = "example"
}

resource "mssql_user" "example" {
  name = mssql_login.example.name
  database = "example-db"
}


resource "mssql_rolemapping" "example" {
  user = mssql_user.example.name
  database = mssql_user.example.database
  role = "db_datareader"
}