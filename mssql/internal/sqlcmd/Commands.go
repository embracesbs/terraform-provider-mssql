package sqlcmd

import (
	"context"
	"database/sql"
)

type ISqlCommand interface {
	Init(connectionString string) error

	Execute(command string, args ...interface{}) error
	Query(query string, scanner func(*sql.Rows) error, args ...interface{}) (*sql.Rows, error)
	Test() string
}

type SqlCommand struct {
	sqlClient  *sql.DB
	connection string
}

func (c *SqlCommand) Init(connectionString string) error {
	var err error
	c.connection = connectionString
	c.sqlClient, err = sql.Open("sqlserver", connectionString)
	if err != nil {
		return err
	}
	return nil
}

func (c *SqlCommand) Execute(command string, args ...interface{}) error {
	ctx := context.Background()

	defer c.sqlClient.Close()

	_, err := c.sqlClient.ExecContext(ctx, command, args...)

	return err
}

func (c *SqlCommand) Query(query string, scanner func(*sql.Rows) error, args ...interface{}) (*sql.Rows, error) {
	ctx := context.Background()

	defer c.sqlClient.Close()

	rows, err := c.sqlClient.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = scanner(rows)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (c *SqlCommand) Test() string {
	return c.connection
}
