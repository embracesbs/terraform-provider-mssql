package mssql

import (
	"context"
	"fmt"
	"log"
	"strconv"

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
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
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
			"owner_sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_read_only": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceDatabaseRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	id := data.Get("id").(string)

	client := meta.(*sqlcmd.SqlCommand)

	database, err := client.GetDatabase(id, "database_id")

	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("name", database.Name)
	data.Set("collation", database.Collation_name)
	data.Set("recovery_mode", database.RecoveryModel)
	data.Set("owner_sid", database.Ownersid)
	data.Set("is_read_only", database.IsReadOnly)

	return diags
}

func resourceDatabaseCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	name := data.Get("name").(string)
	collation := data.Get("collation").(string)
	recovery := data.Get("recovery_mode").(string)

	client := meta.(*sqlcmd.SqlCommand)

	db, err := client.CreateDatabase(name, collation, recovery)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(db.Id))
	data.Set("collation", db.Collation_name)
	data.Set("recovery_mode", db.RecoveryModel)
	data.Set("owner_sid", db.Ownersid)
	data.Set("is_read_only", db.IsReadOnly)

	return diags
}

func resourceDatabaseUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	if data.HasChange("recovery_mode") {

		recovery := data.Get("recovery_mode").(string)
		database := data.Get("name").(string)

		client := meta.(*sqlcmd.SqlCommand)

		err := client.SetRecoveryMode(database, recovery)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	return diags

}

func resourceDatabaseDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	name := data.Get("name").(string)

	client := meta.(*sqlcmd.SqlCommand)

	err := client.DeleteDatabase(name)

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
