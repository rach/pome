package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron"
	"strings"
)

func scheduleMetric(c *cron.Cron, schedule string, fct func()) {
	c.AddFunc(schedule, fct)
	if strings.HasPrefix(schedule, "@every") {
		fct()
	}
}

func initScheduler(
	db *sqlx.DB,
	metrics *MetricList,
	scheduleTableBloat string,
	scheduleDbSize string,
	scheduleIndexBloat string,
	scheduleNumConn string,
	scheduleTPS string) {

	c := cron.New()
	scheduleMetric(c, scheduleIndexBloat,
		func() {
			indexBloatUpdate(db, metrics, GetIndexBloatResult, 120)
		},
	)
	scheduleMetric(c, scheduleTableBloat,
		func() {
			tableBloatUpdate(db, metrics, GetTableBloatResult, 120)
		},
	)
	scheduleMetric(c, scheduleDbSize,
		func() {
			databaseSizeUpdate(db, metrics, GetDatabeSizeResult, 120)
		},
	)
	scheduleMetric(c, scheduleNumConn,
		func() {
			numberOfConnectionUpdate(db, metrics, GetNumberOfConnectionResult, 120)
		},
	)
	scheduleMetric(c, scheduleTPS,
		func() {
			// GetTransactionNumberResult only return the current number of transaction currently
			transactionPerSecUpdate(db, metrics, GetTransactionNumberResult, 120)
		},
	)
	c.Start()
}
