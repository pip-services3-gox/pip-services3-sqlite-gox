package test

import (
	"context"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	persist "github.com/pip-services3-gox/pip-services3-sqlite-gox/persistence"
	"github.com/pip-services3-gox/pip-services3-sqlite-gox/test/fixtures"
)

type DummyRefSqlitePersistence struct {
	persist.IdentifiableSqlitePersistence[*fixtures.Dummy, string]
}

func NewDummyRefSqlitePersistence() *DummyRefSqlitePersistence {
	c := &DummyRefSqlitePersistence{}
	c.IdentifiableSqlitePersistence = *persist.InheritIdentifiableSqlitePersistence[*fixtures.Dummy, string](c, "dummies")
	return c
}

func (c *DummyRefSqlitePersistence) GetPageByFilter(ctx context.Context, correlationId string,
	filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[*fixtures.Dummy], err error) {

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

func (c *DummyRefSqlitePersistence) GetCountByFilter(ctx context.Context, correlationId string,
	filter cdata.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	return c.IdentifiableSqlitePersistence.GetCountByFilter(ctx, correlationId, filterObj)
}
