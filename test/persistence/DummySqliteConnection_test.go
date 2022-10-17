package test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	conn "github.com/pip-services3-gox/pip-services3-sqlite-gox/connect"
	tf "github.com/pip-services3-gox/pip-services3-sqlite-gox/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestDummySqliteConnection(t *testing.T) {

	var persistence *DummySqlitePersistence
	var fixture tf.DummyPersistenceFixture
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

	persistence = NewDummySqlitePersistence()
	descr := cref.NewDescriptor("pip-services", "connection", "sqlite", "default", "1.0")
	ref := cref.NewReferencesFromTuples(context.Background(), descr, connection)
	persistence.SetReferences(context.Background(), ref)

	fixture = *tf.NewDummyPersistenceFixture(persistence)

	opnErr := connection.Open(context.Background(), "")
	if opnErr != nil {
		t.Error("Error opened connection", opnErr)
		return
	}
	defer connection.Close(context.Background(), "")

	opnErr = persistence.Open(context.Background(), "")
	if opnErr != nil {
		t.Error("Error opened persistence", opnErr)
		return
	}
	defer persistence.Close(context.Background(), "")

	opnErr = persistence.Clear(context.Background(), "")
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("Connection", func(t *testing.T) {
		assert.NotNil(t, connection.GetConnection())
		assert.NotNil(t, connection.GetDatabaseName())
		assert.NotEqual(t, "", connection.GetDatabaseName())
	})

	t.Run("DummySqliteConnection:CRUD", fixture.TestCrudOperations)

	opnErr = persistence.Clear(context.Background(), "")
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummySqliteConnection:Batch", fixture.TestBatchOperations)

}
