package main

import (
	"github.com/rach/pom/Godeps/_workspace/src/github.com/alecthomas/kingpin"
	"os"
)

//go:generate go-bindata -prefix "static/" -pkg main -o bindata.go static/index.html static/build/...

const (
	Version = "0.1.0"
)

var (
	app  = kingpin.New("pom", "A Postgres Monitoring Tool.")
	host = app.Flag("host", "database server host (default: localhost)").
		Short('h').PlaceHolder("HOSTNAME").Default("localhost").String()
	port = app.Flag("port", "database server port (default: 2345)").
		Short('p').Default("2345").PlaceHolder("PORT").Int()
	password = app.Flag("password", "").
			Short('W').PlaceHolder("PASSWORD").String()
	username = app.Flag("username", "").
			Short('U').PlaceHolder("USERNAME").Required().String()
	database = app.Arg("DBNAME", "").Required().String()
)

func parseCmdLine(args []string) (command string, err error) {
	//this is isolated from the main() function to make it more testable
	app.Version(Version)
	return app.Parse(args)
}

func main() {
	kingpin.MustParse(parseCmdLine(os.Args[1:]))
	var metrics = MetricList{}
	var connstring = connectionString(*host, *database, *username, *password)
	db := connectDB(connstring)
	context := &appContext{db, &metrics}
	go metricScheduler(db, &metrics, indexBloatUpdate, GetIndexBloatResult, 12*60*60, 120)
	go metricScheduler(db, &metrics, tableBloatUpdate, GetTableBloatResult, 12*60*60, 120)
	go metricScheduler(db, &metrics, databaseSizeUpdate, GetDatabeSizeResult, 60*60, 120)
	go metricScheduler(db, &metrics, numberOfConnectionUpdate, GetNumberOfConnectionResult, 5*60, 120)
	initWebServer(context)
}
