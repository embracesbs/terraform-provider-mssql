package sqlcmd

import "fmt"

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

func (sqlClient *SqlCommand) GetDatabase(id string, where string) (*Database, error) {

	database := &Database{}

	sqlClient.UseDefault()

	result, err := sqlClient.Query(fmt.Sprintf("select name,database_id,owner_sid,compatibility_level,collation_name,is_read_only,recovery_model_desc,is_fulltext_enabled FROM sys.databases WHERE %s = '%s'", where, id))

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

func (sqlClient *SqlCommand) CreateDatabase(name string, collation string, recoveryMode string) (*Database, error) {

	sqlClient.UseDefault()

	cmd := fmt.Sprintf("CREATE DATABASE %s COLLATE %s;", name, collation)
	err := sqlClient.Execute(cmd)

	if err == nil {
		err = sqlClient.SetRecoveryMode(name, recoveryMode)
		if err == nil {
			database, err := sqlClient.GetDatabase(name, "name")
			if err == nil {
				return database, nil
			}
		}
	}

	return nil, err
}

func (sqlClient *SqlCommand) DeleteDatabase(name string) error {

	sqlClient.UseDefault()
	err := sqlClient.Execute(fmt.Sprintf("DROP DATABASE %s;", name))

	return err
}

func (sqlClient *SqlCommand) SetRecoveryMode(name string, mode string) error {

	sqlClient.UseDefault()
	err := sqlClient.Execute(fmt.Sprintf("ALTER DATABASE %s SET RECOVERY %s", name, mode))

	return err
}
