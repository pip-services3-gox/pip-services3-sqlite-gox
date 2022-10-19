package test

import (
	"context"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	persist "github.com/pip-services3-gox/pip-services3-sqlite-gox/persistence"
)

type DummyMapSqlitePersistence struct {
	*persist.IdentifiableSqlitePersistence[map[string]any, string]
}

func NewDummyMapSqlitePersistence() *DummyMapSqlitePersistence {
	c := &DummyMapSqlitePersistence{}
	c.IdentifiableSqlitePersistence = persist.InheritIdentifiableSqlitePersistence[map[string]any, string](c, "dummies")
	return c
}

func (c *DummyMapSqlitePersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiableSqlitePersistence.DefineSchema()
	c.EnsureSchema("CREATE TABLE " + c.IdentifiableSqlitePersistence.TableName + " (\"id\" TEXT PRIMARY KEY, \"key\" TEXT, \"content\" TEXT)")
	c.EnsureIndex(c.IdentifiableSqlitePersistence.TableName+"_key", map[string]string{"key": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyMapSqlitePersistence) GetPageByFilter(ctx context.Context, correlationId string,
	filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[map[string]any], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	sorting := ""

	return c.IdentifiableSqlitePersistence.GetPageByFilter(ctx, correlationId,
		filterObj, paging, sorting, "",
	)
}

func (c *DummyMapSqlitePersistence) GetCountByFilter(ctx context.Context, correlationId string,
	filter cdata.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	return c.IdentifiableSqlitePersistence.GetCountByFilter(ctx, correlationId, filterObj)
}
