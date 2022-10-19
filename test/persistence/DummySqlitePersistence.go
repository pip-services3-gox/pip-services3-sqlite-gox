package test

import (
	"context"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	persist "github.com/pip-services3-gox/pip-services3-sqlite-gox/persistence"
	"github.com/pip-services3-gox/pip-services3-sqlite-gox/test/fixtures"
)

type DummySqlitePersistence struct {
	*persist.IdentifiableSqlitePersistence[fixtures.Dummy, string]
}

func NewDummySqlitePersistence() *DummySqlitePersistence {
	c := &DummySqlitePersistence{}
	c.IdentifiableSqlitePersistence = persist.InheritIdentifiableSqlitePersistence[fixtures.Dummy, string](c, "dummies")
	return c
}

func (c *DummySqlitePersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiableSqlitePersistence.DefineSchema()
	// Row name must be in double quotes for properly case!!!
	c.EnsureSchema("CREATE TABLE " + c.QuotedTableName() + " (\"id\" TEXT PRIMARY KEY, \"key\" TEXT, \"content\" TEXT)")
	c.EnsureIndex(c.IdentifiableSqlitePersistence.TableName+"_key", map[string]string{"key": "1"}, map[string]string{"unique": "true"})
}

func (c *DummySqlitePersistence) GetPageByFilter(ctx context.Context, correlationId string,
	filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	sorting := ""

	return c.IdentifiableSqlitePersistence.GetPageByFilter(ctx, correlationId,
		filterObj, paging,
		sorting, "",
	)
}

func (c *DummySqlitePersistence) GetCountByFilter(ctx context.Context, correlationId string,
	filter cdata.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	return c.IdentifiableSqlitePersistence.GetCountByFilter(ctx, correlationId, filterObj)
}

func (c *DummySqlitePersistence) GetOneRandom(ctx context.Context, correlationId string) (item fixtures.Dummy, err error) {
	return c.IdentifiableSqlitePersistence.GetOneRandom(ctx, correlationId, "")
}
