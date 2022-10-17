package test_connect

import (
	"context"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	conn "github.com/pip-services3-gox/pip-services3-sqlite-gox/connect"
	"github.com/stretchr/testify/assert"
)

func TestSqliteConnectionResolverConnectionConfigWithParams(t *testing.T) {

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.database", "../../data/test.db",
	)

	resolver := conn.NewSqliteConnectionResolver()
	resolver.Configure(context.Background(), dbConfig)

	config, err := resolver.Resolve(context.Background(), "")
	assert.Nil(t, err)

	assert.NotNil(t, config)
	assert.Equal(t, "../../data/test.db", config)
}
