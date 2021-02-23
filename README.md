# Terraform Provider MSSQL

Run the following command to build the provider

```shell
go build -o terraform-provider-mssql
```

## Test sample configuration

First, build and install the provider.

```shell
make local-install
```

This will create the provider binary. Copy that binary into the location that Terraform will try to find it. Depending on the OS, version and Terraform version this could differ.

## [Contributing](docs/contributing.md) 