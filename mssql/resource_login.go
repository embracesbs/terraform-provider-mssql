package mssql

import (
	"context"
	"database/sql"

	"github.com/embracesbs/terraform-provider-mssql/mssql/internal/sqlcmd"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLogin() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceLoginRead,
		CreateContext: resourceLoginCreate,
		UpdateContext: resourceLoginUpdate,
		DeleteContext: resourceLoginDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceLoginRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	// TODO: all read stuff

	return diags
}

func resourceLoginCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	name := data.Get("name")
	password := data.Get("password")

	cmd := `DECLARE @sql nvarchar(max)
					SET @sql = 'CREATE LOGIN ' + QuoteName(@username) + ' ' +
										 'WITH PASSWORD = ' + QuoteName(@password, '''')
					EXEC (@sql)`

	client := meta.(*sqlcmd.SqlCommand)

	err := client.Execute(cmd, sql.Named("username", name), sql.Named("password", password))

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("1")

	return diags
}

func resourceLoginUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags

}

func resourceLoginDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags

}
