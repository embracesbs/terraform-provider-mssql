package mssql

import (
	"context"
	"database/sql"

	"github.com/embracesbs/terraform-provider-mssql/mssql/internal/sqlcmd"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type User struct {
	name         string
	principal_id int
	sid          string
	type_desc    string
}

//
func resourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceUserRead,
		CreateContext: resourceUserCreate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceUserCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	name := data.Get("name")
	database := data.Get("database")

	cmd := "CREATE USER [" + name.(string) + "] FOR LOGIN [" + name.(string) + "]"

	client := meta.(*sqlcmd.SqlCommand)

	client.UseDb(database.(string))

	err := client.Execute(cmd, sql.Named("username", name))

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(name.(string))

	return diags
}

func resourceUserDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	name := data.Get("name")
	database := data.Get("database")

	cmd := "DROP USER [" + name.(string) + "]"

	client := meta.(*sqlcmd.SqlCommand)

	client.UseDb(database.(string))

	err := client.Execute(cmd, sql.Named("username", name))

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags

}
