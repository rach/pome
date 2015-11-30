package main

import (
	//  "os/user"
	//	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"math"
	"time"
)

func connectDB(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("select 1")
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func connectionString(host string, username string) string {
	// TODO: escape single quote
	return fmt.Sprintf("dbname='%s' user='%s' password='' sslmode=disable", host, username)
}

type MetricList struct {
	TableBloat           map[string]tableBloatMetric `json:"table_bloat"`
	IndexBloat           map[string]indexBloatMetric `json:"index_bloat"`
	TopBloatIndexRatio   []topBloatRatioMetric       `json:"top_index_bloat"`
	TopBloatTableRatio   []topBloatRatioMetric       `json:"top_table_bloat"`
	TotalTableBloatBytes []totalBloatBytesMetric     `json:"total_table_bloat_bytes"`
	TotalIndexBloatBytes []totalBloatBytesMetric     `json:"total_index_bloat_bytes"`
	DatabaseSize         []Metric                    `json:"database_size"`
	NumberOfConnection   []Metric                    `json:"number_of_connection"`
}

type metricFct func(db *sqlx.DB, metrics *MetricList, limit int)

type Metric interface {
	GetTimestamp() int64
}

type numberConnectionMetric struct {
	Timestamp int64 `json:"timestamp"`
	Count     int   `db:"count" json:"count"`
}

type databaseSizeMetric struct {
	Timestamp  int64   `json:"timestamp"`
	TableSize  int     `json:"table_size"`
	IndexSize  int     `json:"index_size"`
	TotalSize  int     `json:"total_size"`
	IndexRatio float64 `json:"index_ratio"`
}

type topBloatRatioMetric struct {
	Timestamp  int64   `json:"timestamp"`
	BloatRatio float64 `json:"bloat_ratio"`
}

type totalBloatBytesMetric struct {
	Timestamp  int64 `json:"timestamp"`
	BloatBytes int   `json:"bloat_bytes"`
}

type bloatMetric struct {
	Timestamp  int64   `json:"timestamp"`
	BloatBytes int     `json:"bloat_bytes"`
	BloatRatio float64 `json:"bloat_ratio"`
}
type tableBloatMetric struct {
	TableSchema string        `db:"schema" json:"table_schema"`
	TableName   string        `db:"table" json:"table_name"`
	Bloat       []bloatMetric `json:"data"`
}

type indexBloatMetric struct {
	TableSchema string        `db:"schema" json:"table_schema"`
	TableName   string        `db:"table" json:"table_name"`
	IndexName   string        `db:"index" json:"index_name"`
	Bloat       []bloatMetric `json:"data"`
}

func (m bloatMetric) GetTimestamp() int64            { return m.Timestamp }
func (m databaseSizeMetric) GetTimestamp() int64     { return m.Timestamp }
func (m numberConnectionMetric) GetTimestamp() int64 { return m.Timestamp }

func indexBloatUpdate(db *sqlx.DB, metrics *MetricList, limit int) {
	timestamp := time.Now().Unix()
	rows, err := db.Query(IndexBloatSql)
	if err != nil {
		log.Fatal(err)
	}

	var total_bytes int = 0
	var top_bloat float64 = 0

	// iterate over each row
	for rows.Next() {
		var key string
		var schema string
		var table string
		var index string
		var bloat_bytes int
		var bloat_ratio float64
		err := rows.Scan(&key, &schema, &table, &index, &bloat_bytes, &bloat_ratio)
		if err != nil {
			log.Fatal(err)
		}

		total_bytes += bloat_bytes
		top_bloat = math.Max(top_bloat, bloat_ratio)

		if (*metrics).IndexBloat == nil {
			(*metrics).IndexBloat = make(map[string]indexBloatMetric)
		}
		if _, ok := (*metrics).IndexBloat[key]; !ok {
			(*metrics).IndexBloat[key] = indexBloatMetric{
				TableSchema: table,
				TableName:   table,
				IndexName:   index}
		}
		m := bloatMetric{Timestamp: timestamp, BloatBytes: bloat_bytes, BloatRatio: bloat_ratio}
		tmp_metrics := append((*metrics).IndexBloat[key].Bloat, m)
		if len(tmp_metrics) > limit {
			tmp_metrics = tmp_metrics[len(tmp_metrics)-limit:]
		}
		v := (*metrics).IndexBloat[key]
		v.Bloat = tmp_metrics
		(*metrics).IndexBloat[key] = v
	}

	tmp_top := append((*metrics).TopBloatIndexRatio, topBloatRatioMetric{timestamp, top_bloat})
	if len(tmp_top) > limit {
		tmp_top = tmp_top[len(tmp_top)-limit:]
	}
	(*metrics).TopBloatIndexRatio = tmp_top

	tmp_total := append((*metrics).TotalIndexBloatBytes, totalBloatBytesMetric{timestamp, total_bytes})
	if len(tmp_total) > limit {
		tmp_total = tmp_total[len(tmp_total)-limit:]
	}
	(*metrics).TotalIndexBloatBytes = tmp_total
}

func tableBloatUpdate(db *sqlx.DB, metrics *MetricList, limit int) {
	// TODO: bad duplicate, need to look into make some part generic
	timestamp := time.Now().Unix()
	rows, err := db.Query(TableBloatSql)
	if err != nil {
		log.Fatal(err)
	}

	var total_bytes int = 0
	var top_bloat float64 = 0

	// iterate over each row
	for rows.Next() {
		var key string
		var schema string
		var table string
		var bloat_bytes int
		var bloat_ratio float64
		err := rows.Scan(&key, &schema, &table, &bloat_bytes, &bloat_ratio)
		if err != nil {
			log.Fatal(err)
		}

		total_bytes += bloat_bytes
		top_bloat = math.Max(top_bloat, bloat_ratio)

		if (*metrics).TableBloat == nil {
			(*metrics).TableBloat = make(map[string]tableBloatMetric)
		}

		if _, ok := (*metrics).TableBloat[key]; !ok {
			(*metrics).TableBloat[key] = tableBloatMetric{
				TableSchema: table,
				TableName:   table}
		}
		m := bloatMetric{Timestamp: timestamp, BloatBytes: bloat_bytes, BloatRatio: bloat_ratio}
		tmp_metrics := append((*metrics).TableBloat[key].Bloat, m)
		if len(tmp_metrics) > limit {
			tmp_metrics = tmp_metrics[len(tmp_metrics)-limit:]
		}
		v := (*metrics).TableBloat[key]
		v.Bloat = tmp_metrics
		(*metrics).TableBloat[key] = v
	}

	tmp_top := append((*metrics).TopBloatTableRatio, topBloatRatioMetric{timestamp, top_bloat})
	if len(tmp_top) > limit {
		tmp_top = tmp_top[len(tmp_top)-limit:]
	}
	(*metrics).TopBloatTableRatio = tmp_top

	tmp_total := append((*metrics).TotalTableBloatBytes, totalBloatBytesMetric{timestamp, total_bytes})
	if len(tmp_total) > limit {
		tmp_total = tmp_total[len(tmp_total)-limit:]
	}
	(*metrics).TotalTableBloatBytes = tmp_total
}

func databaseSizeUpdate(db *sqlx.DB, metrics *MetricList, limit int) {
	timestamp := time.Now().Unix()
	var table_size int
	var index_size int
	var total_size int
	var index_ratio float64
	row := db.QueryRow(DatabaseSizeSql)
	err := row.Scan(&table_size, &index_size, &total_size, &index_ratio)
	if err != nil {
		log.Fatal(err)
	}
	tmp := append((*metrics).DatabaseSize, databaseSizeMetric{timestamp, table_size, index_size, total_size, index_ratio})
	if len(tmp) > limit {
		tmp = tmp[len(tmp)-limit:]
	}
	(*metrics).DatabaseSize = tmp
}

func numberOfConnectionUpdate(db *sqlx.DB, metrics *MetricList, limit int) {
	timestamp := time.Now().Unix()
	var numconn int
	row := db.QueryRow(NumberOfConnectionSql)
	err := row.Scan(&numconn)
	if err != nil {
		log.Fatal(err)
	}
	tmp := append((*metrics).NumberOfConnection, numberConnectionMetric{timestamp, numconn})
	if len(tmp) > limit {
		tmp = tmp[len(tmp)-limit:]
	}
	(*metrics).NumberOfConnection = tmp
}

func expireMetrics(metrics *[]Metric, expire int) []Metric {
	//Not used right now, as we are keeping the last 120 metrics
	//this could be more efficient to copy all the slice after index which respect the limit
	limit := time.Now().Unix() - int64(expire)
	m := make([]Metric, 0)
	for _, v := range *metrics {
		if v.GetTimestamp() >= limit {
			m = append(m, v)
		}
	}
	return m
}

func metricScheduler(db *sqlx.DB, metrics *MetricList, mfct metricFct, delay int, limit int) {
	// It's important to use a tmp variable and not appending directly on the pointer as
	// it can be concurrently accessed
	for {
		mfct(db, metrics, limit)
		time.Sleep(time.Duration(delay) * time.Second)
	}
}
