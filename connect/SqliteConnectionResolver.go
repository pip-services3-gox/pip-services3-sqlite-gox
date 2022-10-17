package connect

import (
	"context"
	"strings"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	"github.com/pip-services3-gox/pip-services3-components-gox/auth"
	ccon "github.com/pip-services3-gox/pip-services3-components-gox/connect"
)

// SqliteConnectionResolver a helper struct  that resolves SQLite connection and credential parameters,
// validates them and generates a connection URI.
// It is able to process multiple connections to SQLite cluster nodes.
//
//	Configuration parameters
//		- connection(s):
//			- discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//			- host:                        host name or IP address
//			- port:                        port number (default: 8082)
//			- database:                    database name
//			- uri:                         resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//			- username:                    user name
//			- password:                    user password
//	References
//		- *:discovery:*:*:1.0             (optional) IDiscovery services
//		- *:credential-store:*:*:1.0      (optional) Credential stores to resolve credentials
type SqliteConnectionResolver struct {
	//The connections resolver.
	ConnectionResolver ccon.ConnectionResolver
	//The credentials resolver.
	CredentialResolver auth.CredentialResolver
}

// NewSqliteConnectionResolver creates new connection resolver
//
//	Returns: *SqliteConnectionResolver
func NewSqliteConnectionResolver() *SqliteConnectionResolver {
	sqliteCon := SqliteConnectionResolver{}
	sqliteCon.ConnectionResolver = *ccon.NewEmptyConnectionResolver()
	sqliteCon.CredentialResolver = *auth.NewEmptyCredentialResolver()
	return &sqliteCon
}

// Configure is configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams configuration parameters to be set.
func (c *SqliteConnectionResolver) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.CredentialResolver.Configure(ctx, config)
}

// SetReferences is sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context,
//		- references crefer.IReferences references to locate the component dependencies.
func (c *SqliteConnectionResolver) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.ConnectionResolver.SetReferences(ctx, references)
	c.CredentialResolver.SetReferences(ctx, references)
}

func (c *SqliteConnectionResolver) validateConnection(correlationId string, connection *ccon.ConnectionParams) error {
	uri := connection.Uri()
	if uri != "" {
		if !strings.HasPrefix(uri, "file://") {
			return cerr.NewConfigError(correlationId, "WRONG_PROTOCOL", "Connection protocol must be file://")
		}
		return nil
	}

	if database, ok := connection.GetAsNullableString("database"); !ok || database == "" {
		return cerr.NewConfigError(correlationId, "NO_DATABASE", "Connection database is not set")
	}
	return nil
}

func (c *SqliteConnectionResolver) validateConnections(correlationId string, connections []*ccon.ConnectionParams) error {
	if connections == nil || len(connections) == 0 {
		return cerr.NewConfigError(correlationId, "NO_CONNECTION", "Database connection is not set")
	}
	for _, connection := range connections {
		err := c.validateConnection(correlationId, connection)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *SqliteConnectionResolver) composeUri(connections []*ccon.ConnectionParams, credential *auth.CredentialParams) string {
	// If there is a uri or database then return it immediately
	for _, connection := range connections {
		uri := connection.Uri()
		if uri != "" {
			// Removing file://
			return uri[7:]
		}

		database, isFetched := connection.GetAsNullableString("database")
		if isFetched && database != "" {
			return database
		}
	}

	return ""
}

// Resolve method are resolves SQLite connection URI from connection and credential parameters.
//
//	Parameters:
//		- ctx context.Context
//		- correlationId  string (optional) transaction id to trace execution through call chain.
//	Returns: uri string, err error resolved URI and error, if this occured.
func (c *SqliteConnectionResolver) Resolve(ctx context.Context, correlationId string) (uri string, err error) {
	connections, err := c.ConnectionResolver.ResolveAll(correlationId)
	if err != nil {
		return "", err
	}
	//Validate connections
	err = c.validateConnections(correlationId, connections)
	if err != nil {
		return "", err
	}
	credential, _ := c.CredentialResolver.Lookup(ctx, correlationId)
	return c.composeUri(connections, credential), nil
}
