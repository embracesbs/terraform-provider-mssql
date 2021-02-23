
# mssql below 2012 is not supported yet.
provider "mssql" {

  user_name = "sa"
  password = "example"
  url = "example.database.windows.net"
  
}