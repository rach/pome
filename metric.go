package main

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type MetricList struct {
	TableBloat            map[string]Metric `json:"table_bloat"`
	IndexBloat            map[string]Metric `json:"index_bloat"`
	TopBloatIndexRatio    []Metric          `json:"top_index_bloat"`
	TopBloatTableRatio    []Metric          `json:"top_table_bloat"`
	TotalTableBloatBytes  []Metric          `json:"total_table_bloat_bytes"`
	TotalIndexBloatBytes  []Metric          `json:"total_index_bloat_bytes"`
	DatabaseSize          []Metric          `json:"database_size"`
	NumberOfConnection    []Metric          `json:"number_of_connection"`
	Version               string            `json:"version"`
	TransactionPerSec     []Metric          `json:"transaction_per_sec"`
	LastTransactionNumber TransactionNumber `json:"-"`
}

type TransactionNumber struct {
	Timestamp int64
	Count     int64
}

func initMetricList(version string) MetricList {
	return MetricList{
		Version:              version,
		TableBloat:           make(map[string]Metric),
		IndexBloat:           make(map[string]Metric),
		TopBloatIndexRatio:   []Metric{},
		TopBloatTableRatio:   []Metric{},
		TotalTableBloatBytes: []Metric{},
		TotalIndexBloatBytes: []Metric{},
		DatabaseSize:         []Metric{},
		NumberOfConnection:   []Metric{},
		TransactionPerSec:    []Metric{},
	}
}

type metricFct func(db *sqlx.DB, metrics *MetricList, datafct databaseResultFct, limit int)

type Metric interface {
}

type Metrics []Metric

type transactionPerSecMetric struct {
	Timestamp int64 `json:"timestamp"`
	Tps       int   `json:"tps"`
}

type numberConnectionMetric struct {
	Timestamp int64 `json:"timestamp"`
	Count     int   `json:"count"`
}

type databaseSizeMetric struct {
	Timestamp  int64   `json:"timestamp"`
	TableSize  int     `json:"table_size"`
	IndexSize  int     `json:"index_size"`
	TotalSize  int     `json:"total_size"`
	IndexRatio float64 `json:"index_ratio"`
}

type topTableBloatRatioMetric struct {
	Timestamp   int64   `json:"timestamp"`
	BloatRatio  float64 `json:"bloat_ratio"`
	TableSchema string  `json:"table_schema"`
	TableName   string  `json:"table_name"`
}

type topIndexBloatRatioMetric struct {
	Timestamp   int64   `json:"timestamp"`
	BloatRatio  float64 `json:"bloat_ratio"`
	TableSchema string  `json:"table_schema"`
	TableName   string  `json:"table_name"`
	IndexName   string  `json:"index_name"`
}

type totalBloatBytesMetric struct {
	Timestamp  int64 `json:"timestamp"`
	BloatBytes int64 `json:"bloat_bytes"`
}

type bloatMetric struct {
	Timestamp  int64   `json:"timestamp"`
	BloatBytes int64   `json:"bloat_bytes"`
	BloatRatio float64 `json:"bloat_ratio"`
}

type tableBloatMetric struct {
	TableSchema string   `json:"table_schema"`
	TableName   string   `json:"table_name"`
	Bloat       []Metric `json:"data"`
}

type indexBloatMetric struct {
	TableSchema string   `json:"table_schema"`
	TableName   string   `json:"table_name"`
	IndexName   string   `json:"index_name"`
	Bloat       []Metric `json:"data"`
}

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func appendAndFilter(list []Metric, m Metric, limit int) []Metric {
	r := append(list, m)
	if len(r) > limit {
		r = r[len(r)-limit:]
	}
	return r
}

func initMapMetric(key string, vm *map[string]Metric, metric Metric) {
	if *vm == nil {
		*vm = make(map[string]Metric)
	}
	if _, ok := (*vm)[key]; !ok {
		(*vm)[key] = metric
	}
}

func indexBloatUpdate(db *sqlx.DB, metrics *MetricList, datafct databaseResultFct, limit int) {
	timestamp := GetTimestamp()
	results := (datafct(db)).([]IndexBloatDatabaseResult)
	var totalBytes int64 = 0
	var topBloatRatio float64 = 0
	var topBloatRatioMetric topIndexBloatRatioMetric

	// iterate over each row
	for _, v := range results {
		if v.Schema == "information_schema" {
			continue
		}
		totalBytes += v.BloatBytes
		if v.BloatRatio > topBloatRatio {
			topBloatRatio = v.BloatRatio
			topBloatRatioMetric = topIndexBloatRatioMetric{
				timestamp,
				topBloatRatio,
				v.Schema,
				v.Table,
				v.Index,
			}
		}
		initMapMetric(
			v.Key,
			&((*metrics).IndexBloat),
			indexBloatMetric{
				TableSchema: v.Schema,
				TableName:   v.Table,
				IndexName:   v.Index})

		m := bloatMetric{Timestamp: timestamp, BloatBytes: v.BloatBytes, BloatRatio: v.BloatRatio}
		current_val := ((*metrics).IndexBloat[v.Key]).(indexBloatMetric)
		tmp_metrics := appendAndFilter(current_val.Bloat, m, limit)

		current_val.Bloat = tmp_metrics
		(*metrics).IndexBloat[v.Key] = current_val
	}

	(*metrics).TopBloatIndexRatio = appendAndFilter(
		(*metrics).TopBloatIndexRatio,
		topBloatRatioMetric,
		limit,
	)

	tbbm := totalBloatBytesMetric{timestamp, totalBytes}
	(*metrics).TotalIndexBloatBytes = appendAndFilter((*metrics).TotalIndexBloatBytes, tbbm, limit)
}

func tableBloatUpdate(db *sqlx.DB, metrics *MetricList, datafct databaseResultFct, limit int) {
	timestamp := GetTimestamp()
	results := (datafct(db)).([]TableBloatDatabaseResult)
	var total_bytes int64 = 0
	var topBloatRatio float64 = 0
	var topBloatRatioMetric topTableBloatRatioMetric

	// iterate over each row
	for _, v := range results {
		if v.Schema == "information_schema" {
			continue
		}
		total_bytes += v.BloatBytes
		if v.BloatRatio > topBloatRatio {
			topBloatRatio = v.BloatRatio
			topBloatRatioMetric = topTableBloatRatioMetric{
				timestamp,
				topBloatRatio,
				v.Schema,
				v.Table,
			}
		}

		initMapMetric(
			v.Key,
			&((*metrics).TableBloat),
			tableBloatMetric{
				TableSchema: v.Schema,
				TableName:   v.Table})

		m := bloatMetric{Timestamp: timestamp, BloatBytes: v.BloatBytes, BloatRatio: v.BloatRatio}
		current_val := ((*metrics).TableBloat[v.Key]).(tableBloatMetric)
		tmp_metrics := appendAndFilter(current_val.Bloat, m, limit)

		current_val.Bloat = tmp_metrics
		(*metrics).TableBloat[v.Key] = current_val
	}

	(*metrics).TopBloatTableRatio = appendAndFilter(
		(*metrics).TopBloatTableRatio,
		topBloatRatioMetric,
		limit,
	)

	tbbm := totalBloatBytesMetric{timestamp, total_bytes}
	(*metrics).TotalTableBloatBytes = appendAndFilter((*metrics).TotalTableBloatBytes, tbbm, limit)
}

func databaseSizeUpdate(db *sqlx.DB, metrics *MetricList, datafct databaseResultFct, limit int) {
	timestamp := GetTimestamp()
	res := (datafct(db)).(DatabaseSizeResult)
	met := databaseSizeMetric{timestamp, res.TableSize, res.IndexSize, res.TotalSize, res.IndexRatio}
	(*metrics).DatabaseSize = appendAndFilter((*metrics).DatabaseSize, met, limit)
}

func numberOfConnectionUpdate(db *sqlx.DB, metrics *MetricList, datafct databaseResultFct, limit int) {
	timestamp := GetTimestamp()
	res := (datafct(db)).(NumberOfConnectionResult)
	met := numberConnectionMetric{timestamp, res.Count}
	(*metrics).NumberOfConnection = appendAndFilter((*metrics).NumberOfConnection, met, limit)
}

func transactionPerSecUpdate(db *sqlx.DB, metrics *MetricList, datafct databaseResultFct, limit int) {
	timestamp := GetTimestamp()
	res := (datafct(db)).(TransactionNumberResult)
	if &metrics.LastTransactionNumber != nil {
		tps := (res.Count - metrics.LastTransactionNumber.Count) / (timestamp - metrics.LastTransactionNumber.Timestamp)
		met := transactionPerSecMetric{timestamp, int(tps)}
		(*metrics).TransactionPerSec = appendAndFilter((*metrics).TransactionPerSec, met, limit)
	}
	metrics.LastTransactionNumber = TransactionNumber{Timestamp: timestamp, Count: res.Count}
}
