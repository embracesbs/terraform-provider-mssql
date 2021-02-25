package sqlcmd

type Database struct {
	Id                 int
	Name               string
	Ownersid           string
	CompatibilityLevel string
	Collation_name     string
	IsReadOnly         bool
	RecoveryModel      string
	FullTextEnabled    string
}

func (c *SqlCommand) GetDatabase(id string, where string) (*Database, error) {

	database := &Database{}

	c.UseDefault()

	result, err := c.Query("select name,database_id,owner_sid,compatibility_level,collation_name,is_read_only,recovery_model_desc,is_fulltext_enabled FROM sys.databases WHERE " + where + " = '" + id + "'")

	if err != nil {
		return database, err
	}

	for result.Next() {
		err := result.Scan(&database.Name, &database.Id, &database.Ownersid, &database.CompatibilityLevel, &database.Collation_name, &database.IsReadOnly, &database.RecoveryModel, &database.FullTextEnabled)

		if err != nil {
			return database, err
		}

	}
	return database, nil
}

func (c *SqlCommand) CreateDatabase(name string, collation string, recoveryMode string) (*Database, error) {

	c.UseDefault()

	cmd := "CREATE DATABASE " + name + " COLLATE " + collation + ";"
	err := c.Execute(cmd)

	if err == nil {
		err = c.SetRecoveryMode(name, recoveryMode)
		if err == nil {
			database, err := c.GetDatabase(name, "name")
			if err == nil {
				return database, nil
			}
		}
	}

	return nil, err
}

func (c *SqlCommand) DeleteDatabase(name string) error {

	c.UseDefault()
	err := c.Execute("DROP DATABASE " + name + ";")

	return err
}

func (c *SqlCommand) SetRecoveryMode(name string, mode string) error {

	c.UseDefault()
	err := c.Execute("ALTER DATABASE " + name + " SET RECOVERY " + mode + "")

	return err
}
