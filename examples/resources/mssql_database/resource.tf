resource "mssql_database" "example" {

    name = "example"
    recovery_mode = "FULL"
    collation = "SQL_Latin1_General_CP1_CI_AS"
    
}
