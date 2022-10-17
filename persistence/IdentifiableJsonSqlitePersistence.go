package persistence

import (
	"context"
	"database/sql"
	"encoding/json"

	cconv "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

// IdentifiableJsonSqlitePersistence is an abstract persistence component that stores data in SQLite in JSON or JSONB fields
// and implements a number of CRUD operations over data items with unique ids.
// The data items must implement IIdentifiable interface.
//
// The JSON table has only two fields: id and data.
//
// In basic scenarios child classes shall only override getPageByFilter,
// getListByFilter or deleteByFilter operations with specific filter function.
// All other operations can be used out of the box.
//
// In complex scenarios child classes can implement additional operations by
// accessing c._collection and c._model properties.
//
//	Configuration parameters
//
//		- collection:                  (optional) SQLite collection name
//		- connection(s):
//			- discovery_key:             (optional) a key to retrieve the connection from IDiscovery
//			- host:                      host name or IP address
//			- port:                      port number (default: 27017)
//			- uri:                       resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                 (optional) a key to retrieve the credentials from ICredentialStore
//			- username:                  (optional) user name
//			- password:                  (optional) user password
//		- options:
//			- connect_timeout:      (optional) number of milliseconds to wait before timing out when connecting a new client (default: 0)
//			- idle_timeout:         (optional) number of milliseconds a client must sit idle in the pool and not be checked out (default: 10000)
//			- max_pool_size:        (optional) maximum number of clients the pool should contain (default: 10)
//
//	References
//		- *:logger:*:*:1.0           (optional) ILogger components to pass log messages components to pass log messages
//		- *:discovery:*:*:1.0        (optional) IDiscovery services
//		- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
//
//	Example:
//		type DummyJsonSqlitePersistence struct {
//			*persist.IdentifiableJsonSqlitePersistence[fixtures.Dummy, string]
//		}
//
//		func NewDummyJsonSqlitePersistence() *DummyJsonSqlitePersistence {
//			c := &DummyJsonSqlitePersistence{}
//			c.IdentifiableJsonSqlitePersistence = persist.InheritIdentifiableJsonSqlitePersistence[fixtures.Dummy, string](c, "dummies_json")
//			return c
//		}
//
//		func (c *DummyJsonSqlitePersistence) DefineSchema() {
//			c.ClearSchema()
//			c.IdentifiableJsonSqlitePersistence.DefineSchema()
//			c.EnsureTable("", "")
//			c.EnsureIndex(c.TableName+"_key", map[string]string{"(data->'key')": "1"}, map[string]string{"unique": "true"})
//		}
//
//		func (c *DummyJsonSqlitePersistence) GetPageByFilter(ctx context.Context, correlationId string,
//			filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[fixtures.Dummy], err error) {
//
//			key, ok := filter.GetAsNullableString("Key")
//			filterObj := ""
//			if ok && key != "" {
//				filterObj += "data->key='" + key + "'"
//			}
//
//			return c.IdentifiableJsonSqlitePersistence.GetPageByFilter(ctx, correlationId,
//				filterObj, paging,
//				"", "",
//			)
//		}
//
//		func (c *DummyJsonSqlitePersistence) GetCountByFilter(ctx context.Context, correlationId string,
//			filter cdata.FilterParams) (count int64, err error) {
//
//			filterObj := ""
//			if key, ok := filter.GetAsNullableString("Key"); ok && key != "" {
//				filterObj += "data->key='" + key + "'"
//			}
//
//			return c.IdentifiableJsonSqlitePersistence.GetCountByFilter(ctx, correlationId, filterObj)
//		}
//
//		func (c *DummyJsonSqlitePersistence) GetOneRandom(ctx context.Context, correlationId string) (item fixtures.Dummy, err error) {
//			return c.IdentifiableJsonSqlitePersistence.GetOneRandom(ctx, correlationId, "")
//		}
type IdentifiableJsonSqlitePersistence[T any, K any] struct {
	*IdentifiableSqlitePersistence[T, K]
}

// InheritIdentifiableJsonSqlitePersistence creates a new instance of the persistence component.
//
//	Parameters:
//		- overrides References to override virtual methods
//		- tableName    (optional) a table name.
func InheritIdentifiableJsonSqlitePersistence[T any, K any](overrides ISqlitePersistenceOverrides[T], tableName string) *IdentifiableJsonSqlitePersistence[T, K] {
	c := &IdentifiableJsonSqlitePersistence[T, K]{}
	c.IdentifiableSqlitePersistence = InheritIdentifiableSqlitePersistence[T, K](overrides, tableName)
	return c
}

// EnsureTable Adds DML statement to automatically create JSON table
//
//		Parameters:
//	  - idType type of the id column (default: VARCHAR(32))
//	  - dataType type of the data column (default: JSON)
func (c *IdentifiableJsonSqlitePersistence[T, K]) EnsureTable(idType string, dataType string) {
	if idType == "" {
		idType = "VARCHAR(32)"
	}
	if dataType == "" {
		dataType = "JSON"
	}

	query := "CREATE TABLE IF NOT EXISTS " + c.QuotedTableName() +
		" (id " + idType + " PRIMARY KEY, data " + dataType + ")"
	c.EnsureSchema(query)
}

// ConvertToPublic converts object value from internal to public format.
//
//	Parameters:
//		- value an object in internal format to convert.
//	Returns: converted object in public format.
func (c *IdentifiableJsonSqlitePersistence[T, K]) ConvertToPublic(rows *sql.Rows) (T, error) {
	var defaultValue T

	columns, err := rows.Columns()
	if err != nil || columns == nil || len(columns) == 0 {
		return defaultValue, err
	}

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}

	err = rows.Scan(pointers...)
	if err != nil {
		return defaultValue, err
	}

	buf := make(map[string]interface{}, 0)

	for index, column := range columns {
		buf[column] = values[index]
	}

	item, ok := buf["data"]
	if !ok {
		item = buf
	}

	_buf, toJsonErr := cconv.JsonConverter.ToJson(item)
	if toJsonErr != nil {
		return defaultValue, toJsonErr
	}

	_item, fromJsonErr := c.IdentifiableSqlitePersistence.JsonConvertor.FromJson(_buf)
	return _item, fromJsonErr
}

// ConvertFromPublic convert object value from public to internal format.
//
//		Parameters:
//	   - value     an object in public format to convert.
//
// Returns converted object in internal format.
func (c *IdentifiableJsonSqlitePersistence[T, K]) ConvertFromPublic(value T) (map[string]any, error) {
	id := GetObjectId[K](value)
	json, _ := json.Marshal(value)
	result := map[string]any{
		"id":   id,
		"data": (string)(json),
	}
	return result, nil
}

// ConvertFromPublicPartial convert object value from public to internal format.
//
//	Parameters:
//		- value     an object in public format to convert.
//	Returns: converted object in internal format.
func (c *IdentifiableJsonSqlitePersistence[T, K]) ConvertFromPublicPartial(value map[string]any) (map[string]any, error) {
	buf, toJsonErr := cconv.JsonConverter.ToJson(value)
	if toJsonErr != nil {
		return nil, toJsonErr
	}
	item, fromJsonErr := c.IdentifiableSqlitePersistence.JsonConvertor.FromJson(buf)
	if toJsonErr != nil {
		return nil, fromJsonErr
	}
	return c.ConvertFromPublic(item)
}

// UpdatePartially updates only few selected fields in a data item.
//
//	Parameters:
//		- ctx context.Context
//		- correlation_id    (optional) transaction id to trace execution through call chain.
//		- id                an id of data item to be updated.
//		- data              a map with fields to be updated.
//
// Returns: receives updated item or error.
func (c *IdentifiableJsonSqlitePersistence[T, K]) UpdatePartially(ctx context.Context, correlationId string,
	id K, data cdata.AnyValueMap) (result T, err error) {

	query := "UPDATE " + c.QuotedTableName() + " SET data=JSON_PATCH(data,?) WHERE id=?"
	values := []any{id, data.Value()}

	qResult, err := c.IdentifiableSqlitePersistence.Client.QueryContext(ctx, query, values...)
	if err != nil {
		return result, err
	}
	defer qResult.Close()

	if !qResult.Next() {
		return result, qResult.Err()
	}

	result, convErr := c.Overrides.ConvertToPublic(qResult)
	if convErr == nil {
		c.IdentifiableSqlitePersistence.Logger.Trace(ctx, correlationId, "Updated partially in %s with id = %s", c.IdentifiableSqlitePersistence.TableName, id)
		return result, nil
	} else {
		return result, convErr
	}
}
