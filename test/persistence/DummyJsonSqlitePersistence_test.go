package test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	tf "github.com/pip-services3-gox/pip-services3-sqlite-gox/test/fixtures"
)

func TestDummyJsonSqlitePersistence(t *testing.T) {

	var persistence *DummyJsonSqlitePersistence
	var fixture tf.DummyPersistenceFixture

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

	persistence = NewDummyJsonSqlitePersistence()
	fixture = *tf.NewDummyPersistenceFixture(persistence)
	persistence.Configure(context.Background(), dbConfig)

	opnErr := persistence.Open(context.Background(), "")
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

	t.Run("DummySqliteConnection:CRUD", fixture.TestCrudOperations)

	opnErr = persistence.Clear(context.Background(), "")
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummySqliteConnection:Batch", fixture.TestBatchOperations)

}
