package build

import (
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cbuild "github.com/pip-services3-gox/pip-services3-components-gox/build"
	conn "github.com/pip-services3-gox/pip-services3-sqlite-gox/connect"
)

// DefaultSqliteFactory helps creates Sqlite components by their descriptors.
//
//	see Factory
//	see SqliteConnection
type DefaultSqliteFactory struct {
	cbuild.Factory
}

// NewDefaultSqliteFactory are create a new instance of the factory.
//
//	Returns: *DefaultSqliteFactory
func NewDefaultSqliteFactory() *DefaultSqliteFactory {
	c := DefaultSqliteFactory{}

	sqliteConnectionDescriptor := cref.NewDescriptor("pip-services", "connection", "sqlite", "*", "1.0")

	c.RegisterType(sqliteConnectionDescriptor, conn.NewSqliteConnection)
	return &c
}
