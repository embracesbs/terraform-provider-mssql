package mssql

import (
	"context"
	"database/sql"

	"github.com/embracesbs/terraform-provider-mssql/mssql/internal/sqlcmd"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Login struct {
	name         string
	principal_id int
	sid          string
	type_desc    string
}

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
			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceLoginRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

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

	client.UseDefault()

	err := client.Execute(cmd, sql.Named("username", name), sql.Named("password", password))

	if err != nil {
		return diag.FromErr(err)
	}

	createdUser, err := getSqlUserDetails(name.(string), client)

	if err != nil {
		return diag.FromErr(err)
	}

	if createdUser != nil {
		data.Set("sid", createdUser.sid)
		data.Set("login_type", createdUser.type_desc)

	}

	data.SetId(name.(string))

	return diags
}

func resourceLoginUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	name := data.Get("name")

	if data.HasChange("password") {

		password := data.Get("password")

		client := meta.(*sqlcmd.SqlCommand)

		client.UseDefault()

		err := client.Execute("ALTER LOGIN ["+name.(string)+"] WITH PASSWORD = '"+password.(string)+"'", sql.Named("username", name))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	return diags

}

func resourceLoginDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	name := data.Get("name")

	client := meta.(*sqlcmd.SqlCommand)

	client.UseDefault()

	err := client.Execute("DROP LOGIN ["+name.(string)+"]", sql.Named("username", name))

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags

}

//Refacter to some sort of generic function that we can use on all resoucres to convert them to struct
func getSqlUserDetails(userName string, client *sqlcmd.SqlCommand) (*Login, error) {

	var logins []*Login

	rows, err := client.Query("SELECT name,principal_id,sid,type_desc from master.sys.sql_logins WHERE [name] = @username", sql.Named("username", userName))

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := new(Login)
		err := rows.Scan(&c.name, &c.principal_id, &c.sid, &c.type_desc)
		if err != nil {
			return nil, err
		}

		logins = append(logins, c) // add each instance to the slice
	}
	if err := rows.Err(); err != nil { // make sure that there was no issue during the process
		return nil, err
	}

	if len(logins) < 1 {
		return nil, nil
	}

	return logins[0], nil

}
