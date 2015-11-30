package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/elazarl/go-bindata-assetfs"
	"net/http"
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
	db, _ := connectDB(connstring)
	context := &appContext{db, &metrics}
	go metricScheduler(db, &metrics, indexBloatUpdate, 12*60*60, 120)
	go metricScheduler(db, &metrics, tableBloatUpdate, 12*60*60, 120)
	go metricScheduler(db, &metrics, databaseSizeUpdate, 60*60, 120)
	go metricScheduler(db, &metrics, numberOfConnectionUpdate, 5*60, 120)
	http.Handle("/api/stats", appHandler{context, metricsHandler})
	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
