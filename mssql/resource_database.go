package mssql

import (
	"context"
	"fmt"
	"log"

	"github.com/embracesbs/terraform-provider-mssql/mssql/internal/sqlcmd"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDatabaseRead,
		CreateContext: resourceDatabaseCreate,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"collation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "SQL_Latin1_General_CP1_CI_AS",
			},
			"recovery_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SIMPLE",
				ValidateFunc: validateRecoveryModel,
			},
		},
	}
}

func resourceDatabaseRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceDatabaseCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	database := data.Get("name")
	collation := data.Get("collation")
	recovery := data.Get("recovery_mode")

	cmd := "CREATE DATABASE " + database.(string) + " COLLATE " + collation.(string) + "; ALTER DATABASE " + database.(string) + " SET RECOVERY " + recovery.(string) + ""

	client := meta.(*sqlcmd.SqlCommand)

	client.UseDefault()

	err := client.Execute(cmd)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(database.(string))

	return diags
}

func resourceDatabaseUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	if data.HasChange("recovery_mode") {

		recovery := data.Get("recovery_mode")
		database := data.Get("name")

		client := meta.(*sqlcmd.SqlCommand)

		client.UseDefault()

		err := client.Execute("ALTER DATABASE " + database.(string) + " SET RECOVERY " + recovery.(string) + "")

		if err != nil {
			return diag.FromErr(err)
		}

	}

	return diags

}

func resourceDatabaseDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	database := data.Get("name")

	cmd := "DROP DATABASE " + database.(string) + ";"

	client := meta.(*sqlcmd.SqlCommand)

	client.UseDefault()

	err := client.Execute(cmd)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags

}

func validateRecoveryModel(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "SIMPLE" || value == "FULL" {
		log.Println("correct recovery model")
	} else {
		errors = append(errors, fmt.Errorf("%g is not valid, only 'FULL' and 'SIMPLE'", value))
	}

	return ws, errors

}
