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
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mssql_login":       resourceLogin(),
			"mssql_user":        resourceUser(),
			"mssql_rolemapping": resourceRoleMapping(),
			"mssql_database":    resourceDatabase(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

var MssqlClient sqlcmd.ISqlCommand

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	MssqlClient = &sqlcmd.SqlCommand{}

	var diags diag.Diagnostics
	user := d.Get("user_name").(string)
	password := d.Get("password").(string)
	url := d.Get("url").(string)

	err := MssqlClient.Init(user, password, url)

	if err != nil {

		return nil, diag.FromErr(err)
	}

	return MssqlClient, diags

}
