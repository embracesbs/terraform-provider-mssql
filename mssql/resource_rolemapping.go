package mssql

import (
	"context"
	"database/sql"

	"github.com/embracesbs/terraform-provider-mssql/mssql/internal/sqlcmd"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRoleMapping() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRoleMappingRead,
		CreateContext: resourceRoleMappingCreate,
		DeleteContext: resourceRoleMappingDelete,
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRoleMappingRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceRoleMappingCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	user := data.Get("user")
	database := data.Get("database")
	role := data.Get("role")

	cmd := "ALTER ROLE [" + role.(string) + "] ADD MEMBER " + user.(string) + ""

	client := meta.(*sqlcmd.SqlCommand)

	client.UseDb(database.(string))

	err := client.Execute(cmd, sql.Named("username", user))

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(user.(string) + database.(string) + role.(string))

	return diags
}

func resourceRoleMappingDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	user := data.Get("user")
	database := data.Get("database")
	role := data.Get("role")

	cmd := "ALTER ROLE [" + role.(string) + "] DROP MEMBER " + user.(string) + ""

	client := meta.(*sqlcmd.SqlCommand)

	client.UseDb(database.(string))

	err := client.Execute(cmd, sql.Named("username", user))

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags

}
