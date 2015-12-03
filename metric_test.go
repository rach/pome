package main

import (
	"github.com/rach/poda/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	"testing"
)

var IndexBloatSqlStub = []IndexBloatDatabaseResult{
	{"public.foo.bar", "public", "foo", "bar", 1000, 10.},
	{"public.foo.pk", "public", "foo", "pk", 2000, 20.},
}

var TableBloatSqlStub = []TableBloatDatabaseResult{
	{"public.foo", "public", "foo", 1000, 10.},
	{"public.bar", "public", "bar", 2000, 20.},
}

var DatabaseSizeSqlStub = DatabaseSizeResult{10000, 50000, 15000, 50.}
var NumberOfConnectionSqlStub = NumberOfConnectionResult{5}

func getIndexBloatStub(db *sqlx.DB) interface{} {
	return IndexBloatSqlStub
}

func getTableBloatStub(db *sqlx.DB) interface{} {
	return TableBloatSqlStub
}

func getDatabaseSizeSqlStub(db *sqlx.DB) interface{} {
	return DatabaseSizeSqlStub
}

func getNumberOfConnectionSqlStub(db *sqlx.DB) interface{} {
	return NumberOfConnectionSqlStub
}

func TestIndexBloatUpdate(t *testing.T) {
	metrics := &MetricList{}
	db := &sqlx.DB{} // we are doing real query so don't connect

	indexBloatUpdate(db, metrics, getIndexBloatStub, 1)
	if len((*metrics).IndexBloat) != 2 {
		// number of keys in the map public.foo.bar & public.foo.pk
		t.Error("Index bloat map excepted equal to 2")
	}

	if len((*metrics).TopBloatIndexRatio) != 1 {
		t.Error("Top Bloat Index Ratio excepted equal to 1")
	}

	if len((*metrics).TotalIndexBloatBytes) != 1 {
		t.Error("Total Index Bloat Bytes excepted equal to 1")
	}
	// Testing if the keys has been created in the map
	if _, ok := (*metrics).IndexBloat[IndexBloatSqlStub[0].Key]; !ok {
		t.Error("Index bloat map excepted key public.foo.bar")
	}

	if _, ok := (*metrics).IndexBloat[IndexBloatSqlStub[1].Key]; !ok {
		t.Error("Index bloat map excepted key public.foo.pk")
	}
	// Testing the limit
	indexBloatUpdate(db, metrics, getIndexBloatStub, 1)

	if len((*metrics).IndexBloat) != 2 {
		// number of keys in the map public.foo.bar & public.foo.pk
		t.Error("Index bloat map still excepted equal to 2")
	}

	if len((*metrics).TopBloatIndexRatio) != 1 {
		t.Error("Top Bloat Index Ratio excepted equal to 1")
	}

	if len((*metrics).TotalIndexBloatBytes) != 1 {
		t.Error("Total Index Bloat Bytes excepted still equal to 1")
	}

	tibb := IndexBloatSqlStub[0].BloatBytes + IndexBloatSqlStub[1].BloatBytes
	// Testing the Total
	if len((*metrics).TotalIndexBloatBytes) != 1 && (*metrics).TotalIndexBloatBytes[0] == tibb {
		t.Error("Total Index Bloat Bytes excepted still equal to 1")
	}

	var tibr float64
	if IndexBloatSqlStub[0].BloatRatio > IndexBloatSqlStub[1].BloatRatio {
		tibr = IndexBloatSqlStub[0].BloatRatio
	} else {
		tibr = IndexBloatSqlStub[1].BloatRatio
	}

	if len((*metrics).TopBloatIndexRatio) != 1 && (*metrics).TopBloatIndexRatio[0] == tibr {
		t.Error("Total Index Bloat Bytes excepted still equal to 1")
	}

}

func TestTableBloatUpdate(t *testing.T) {
	metrics := &MetricList{}
	db := &sqlx.DB{} // we are doing real query so don't connect

	tableBloatUpdate(db, metrics, getTableBloatStub, 1)
	if len((*metrics).TableBloat) != 2 {
		// number of keys in the map public.foo.bar & public.foo.pk
		t.Error("Table bloat map excepted equal to 2")
	}

	if len((*metrics).TopBloatTableRatio) != 1 {
		t.Error("Top Bloat Table Ratio excepted equal to 1")
	}

	if len((*metrics).TotalTableBloatBytes) != 1 {
		t.Error("Total Table Bloat Bytes excepted equal to 1")
	}
	// Testing if the keys has been created in the map
	if _, ok := (*metrics).TableBloat[TableBloatSqlStub[0].Key]; !ok {
		t.Error("Table bloat map excepted key public.foo.bar")
	}

	if _, ok := (*metrics).TableBloat[TableBloatSqlStub[1].Key]; !ok {
		t.Error("Table bloat map excepted key public.foo.pk")
	}
	// Testing the limit
	tableBloatUpdate(db, metrics, getTableBloatStub, 1)

	if len((*metrics).TableBloat) != 2 {
		// number of keys in the map public.foo.bar & public.foo.pk
		t.Error("Table bloat map still excepted equal to 2")
	}

	if len((*metrics).TopBloatTableRatio) != 1 {
		t.Error("Top Bloat Table Ratio excepted equal to 1")
	}

	if len((*metrics).TotalTableBloatBytes) != 1 {
		t.Error("Total Table Bloat Bytes excepted still equal to 1")
	}

	tibb := TableBloatSqlStub[0].BloatBytes + TableBloatSqlStub[1].BloatBytes
	// Testing the Total
	if len((*metrics).TotalTableBloatBytes) != 1 && (*metrics).TotalTableBloatBytes[0] == tibb {
		t.Error("Total Table Bloat Bytes excepted still equal to 1")
	}

	var tibr float64
	if TableBloatSqlStub[0].BloatRatio > TableBloatSqlStub[1].BloatRatio {
		tibr = TableBloatSqlStub[0].BloatRatio
	} else {
		tibr = TableBloatSqlStub[1].BloatRatio
	}

	if len((*metrics).TopBloatTableRatio) != 1 && (*metrics).TopBloatTableRatio[0] == tibr {
		t.Error("Total Table Bloat Bytes excepted still equal to 1")
	}
}

func TestDatabaseSizeUpdate(t *testing.T) {
	metrics := &MetricList{}
	db := &sqlx.DB{} // we are doing real query so we don't care

	databaseSizeUpdate(db, metrics, getDatabaseSizeSqlStub, 1)

	length := len((*metrics).DatabaseSize)
	if length != 1 {
		t.Error("DatabaseSize length excepted equal to 1")
	}

	// Testing the limit
	databaseSizeUpdate(db, metrics, getDatabaseSizeSqlStub, 1)

	length = len((*metrics).DatabaseSize)
	if length != 1 {
		t.Error("DatabaseSize length excepted equal to 1")
	}

	// Testing the Values
	if length >= 1 {
		dsm := ((*metrics).DatabaseSize[0]).(databaseSizeMetric)

		if dsm.TableSize != DatabaseSizeSqlStub.TableSize {
			t.Error("wrong value for table size")
		}
		if dsm.TotalSize != DatabaseSizeSqlStub.TotalSize {
			t.Error("wrong value for total size")
		}
		if dsm.IndexSize != DatabaseSizeSqlStub.IndexSize {
			t.Error("wrong value for index size")
		}
		if dsm.IndexRatio != DatabaseSizeSqlStub.IndexRatio {
			t.Error("wrong value for index ratio")
		}
	}

}

func TestNumberOfConnectionUpdate(t *testing.T) {
	metrics := &MetricList{}
	db := &sqlx.DB{} // we are doing real query so we don't care

	numberOfConnectionUpdate(db, metrics, getNumberOfConnectionSqlStub, 1)

	length := len((*metrics).NumberOfConnection)
	if length != 1 {
		t.Error("NumberOfConnection length excepted equal to 1")
	}

	// Testing the limit
	numberOfConnectionUpdate(db, metrics, getNumberOfConnectionSqlStub, 1)

	length = len((*metrics).NumberOfConnection)
	if length != 1 {
		t.Error("NumberOfConnection length excepted equal to 1")
	}

	// Testing the Values
	if length >= 1 {
		noc := ((*metrics).NumberOfConnection[0]).(numberConnectionMetric)

		if noc.Count != NumberOfConnectionSqlStub.Count {
			t.Error("wrong value for count")
		}
	}
}

func TestAppendAndFilter(t *testing.T) {
	l := []Metric{}
	l = appendAndFilter(l, databaseSizeMetric{}, 1)
	if len(l) != 1 {
		t.Error("length expected to be 1")
	}
	l = appendAndFilter(l, databaseSizeMetric{}, 1)
	if len(l) != 1 {
		t.Error("length expected to still be 1 because of the limit")
	}
}
