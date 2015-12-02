package main

import (
	"github.com/rach/pomod/Godeps/_workspace/src/github.com/alecthomas/kingpin"
)

//go:generate go-bindata -prefix "static/" -pkg main -o bindata.go static/index.html static/build/...

var (
	host = kingpin.Flag("host", "database server host (default: localhost)").
		Short('h').PlaceHolder("HOSTNAME").Default("localhost").String()
	port = kingpin.Flag("port", "database server port (default: 5432)").
		Short('p').Default("2345").PlaceHolder("PORT").Int()
	username = kingpin.Flag("username", "").
			Short('U').PlaceHolder("USERNAME").Required().String()
	database = kingpin.Arg("DBNAME", "").Required().String()
)

func main() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
	var metrics = MetricList{}
	var connstring = connectionString(*host, *username)
	db := connectDB(connstring)
	context := &appContext{db, &metrics}
	go metricScheduler(db, &metrics, indexBloatUpdate, GetIndexBloatResult, 12*60*60, 120)
	go metricScheduler(db, &metrics, tableBloatUpdate, GetTableBloatResult, 12*60*60, 120)
	go metricScheduler(db, &metrics, databaseSizeUpdate, GetDatabeSizeResult, 60*60, 120)
	go metricScheduler(db, &metrics, numberOfConnectionUpdate, GetNumberOfConnectionResult, 5*60, 120)
	initWebServer(context)
}
