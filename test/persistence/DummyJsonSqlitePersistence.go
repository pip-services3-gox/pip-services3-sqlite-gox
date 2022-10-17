package test

import (
	"context"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	persist "github.com/pip-services3-gox/pip-services3-sqlite-gox/persistence"
	"github.com/pip-services3-gox/pip-services3-sqlite-gox/test/fixtures"
)

type DummyJsonSqlitePersistence struct {
	*persist.IdentifiableJsonSqlitePersistence[fixtures.Dummy, string]
}

func NewDummyJsonSqlitePersistence() *DummyJsonSqlitePersistence {
	c := &DummyJsonSqlitePersistence{}
	c.IdentifiableJsonSqlitePersistence = persist.InheritIdentifiableJsonSqlitePersistence[fixtures.Dummy, string](c, "dummies_json")
	return c
}

func (c *DummyJsonSqlitePersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiableJsonSqlitePersistence.DefineSchema()
	c.EnsureTable("", "")
	c.EnsureIndex(c.TableName+"_json_key", map[string]string{"(data->'key')": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyJsonSqlitePersistence) GetPageByFilter(ctx context.Context, correlationId string,
	filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "data->key='" + key + "'"
	}

	return c.IdentifiableJsonSqlitePersistence.GetPageByFilter(ctx, correlationId,
		filterObj, paging,
		"", "",
	)
}

func (c *DummyJsonSqlitePersistence) GetCountByFilter(ctx context.Context, correlationId string,
	filter cdata.FilterParams) (count int64, err error) {

	filterObj := ""
	if key, ok := filter.GetAsNullableString("Key"); ok && key != "" {
		filterObj += "data->key='" + key + "'"
	}

	return c.IdentifiableJsonSqlitePersistence.GetCountByFilter(ctx, correlationId, filterObj)
}

func (c *DummyJsonSqlitePersistence) GetOneRandom(ctx context.Context, correlationId string) (item fixtures.Dummy, err error) {
	return c.IdentifiableJsonSqlitePersistence.GetOneRandom(ctx, correlationId, "")
}
