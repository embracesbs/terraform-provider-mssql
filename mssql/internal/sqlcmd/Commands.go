package sqlcmd

import (
	"context"
	"database/sql"
)

type ISqlCommand interface {
	Init(username string, password string, url string) error

	UseDefault() error

	Execute(command string, args ...interface{}) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type SqlCommand struct {
	sqlClient  *sql.DB
	userName   string
	password   string
	url        string
	connection string
}

func (c *SqlCommand) Init(username string, password string, url string) error {
	var err error

	c.userName = username
	c.password = password
	c.url = url

	connectionString := "Server=" + c.url + ";Persist Security Info=False;User ID=" + c.userName + ";Password=" + c.password + ";"
	c.connection = connectionString
	c.sqlClient, err = sql.Open("sqlserver", connectionString)
	if err != nil {
		return err
	}
	return nil
}

func (c *SqlCommand) UseDefault() error {
	var err error

	c.sqlClient, err = sql.Open("sqlserver", c.connection)
	if err != nil {
		return err
	}
	return nil
}

func (c *SqlCommand) UseDb(databaseName string) error {
	var err error

	connectionString := "Server=" + c.url + ";database=" + databaseName + "; Persist Security Info=False;User ID=" + c.userName + ";Password=" + c.password + ";"
	c.sqlClient, err = sql.Open("sqlserver", connectionString)
	if err != nil {
		return err
	}
	return nil
}

func (c *SqlCommand) Execute(command string, args ...interface{}) error {
	ctx := context.Background()

	_, err := c.sqlClient.ExecContext(ctx, command, args...)

	return err
}

func (c *SqlCommand) Query(query string, args ...interface{}) (*sql.Rows, error) {
	ctx := context.Background()

	rows, err := c.sqlClient.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
