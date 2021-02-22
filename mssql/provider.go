package mssql

import (
	"context"

	"github.com/embracesbs/terraform-provider-mssql/mssql/internal/sqlcmd"

	// import sqlserver driver
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"connection_string": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap:         map[string]*schema.Resource{},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

var MssqlClient sqlcmd.ISqlCommand

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	connection_string := d.Get("connection_string").(string)

	err := MssqlClient.Init(connection_string)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return MssqlClient, diags

}