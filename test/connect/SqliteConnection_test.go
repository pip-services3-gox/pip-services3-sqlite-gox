package test_connect

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	conn "github.com/pip-services3-gox/pip-services3-sqlite-gox/connect"
	"github.com/stretchr/testify/assert"
)

func TestSqliteConnection(t *testing.T) {
	var connection *conn.SqliteConnection

	sqliteDatabase := os.Getenv("SQLITE_DB")
	if sqliteDatabase == "" {
		sqliteDatabase = "../../data/test.db"
	}

	if sqliteDatabase == "" {
		panic("Connection params losse")
	}

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.database", sqliteDatabase,
	)

	connection = conn.NewSqliteConnection()
	connection.Configure(context.Background(), dbConfig)
	err := connection.Open(context.Background(), "")
	assert.Nil(t, err)

	defer connection.Close(context.Background(), "")

	assert.NotNil(t, connection.GetConnection())
	assert.NotNil(t, connection.GetDatabaseName())
	assert.NotEqual(t, "", connection.GetDatabaseName())
}
