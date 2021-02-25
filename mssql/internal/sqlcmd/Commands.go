package sqlcmd

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
)

const SqlServer2019 int = 15
const SqlServer2017 int = 14
const SqlServer2016 int = 13
const SqlServer2014 int = 12
const SqlServer2012 int = 11
const SqlServer2008 int = 10
const Sqlserver2005 int = 9

type ISqlCommand interface {
	Init(username string, password string, url string) error

	UseDefault() error

	Execute(command string, args ...interface{}) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	GetVersion() int
}

type SqlCommand struct {
	sqlClient   *sql.DB
	userName    string
	password    string
	url         string
	connection  string
	mainVersion int
	version     string
	edition     string
}

func (c *SqlCommand) GetVersion() int {
	return c.mainVersion
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

	c.SetEdition()

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

func (c *SqlCommand) SetEdition() {

	result, _ := c.Query("SELECT SERVERPROPERTY('productversion') as  'version', SERVERPROPERTY ('edition') as 'edition'")

	for result.Next() {
		result.Scan(&c.version, &c.edition)
	}

	c.mainVersion, _ = strconv.Atoi(strings.Split(c.version, ".")[0])

}
