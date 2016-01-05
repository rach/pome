package main

import (
	"fmt"
	"github.com/rach/pome/Godeps/_workspace/src/github.com/alecthomas/kingpin"
	"log"
	"os"
	"os/user"
)

//go:generate go-bindata -prefix "static/" -pkg main -o bindata.go static/index.html static/build/...

const (
	Version = "0.1.1"
)

func addUsernameFlag(app *kingpin.Application) *string {
	u, err := user.Current()
	if err != nil {
		return app.Flag("username", "").Short('U').
			PlaceHolder("USERNAME").Required().String()
	}
	return app.Flag("username", "").Short('U').Default(u.Username).
		PlaceHolder(fmt.Sprintf("USERNAME (default: %s)", u.Username)).String()
}

var (
	app  = kingpin.New("pome", "A Postgres Metrics Dashboard.")
	host = app.Flag("host", "database server host (default: localhost)").
		Short('h').PlaceHolder("HOSTNAME").Default("localhost").String()
	port = app.Flag("port", "database server port (default: 2345)").
		Short('p').Default("2345").PlaceHolder("PORT").Int()
	password = app.Flag("password", "").Short('W').Bool()
	username = addUsernameFlag(app)
	database = app.Arg("DBNAME", "").Required().String()
)

func parseCmdLine(args []string) (command string, err error) {
	//this is isolated from the main() function to make it more testable
	app.Version(Version)
	return app.Parse(args)
}

func main() {
	kingpin.MustParse(parseCmdLine(os.Args[1:]))
	var metrics = MetricList{Version: Version}
	pwd := ""
	if *password {
		fmt.Print("Enter Password: ")
		fmt.Scanln(&pwd)
	}
	var connstring = connectionString(*host, *database, *username, pwd)
	db := connectDB(connstring)
	context := &appContext{db, &metrics}
	go metricScheduler(db, &metrics, indexBloatUpdate, GetIndexBloatResult, 12*60*60, 120)
	go metricScheduler(db, &metrics, tableBloatUpdate, GetTableBloatResult, 12*60*60, 120)
	go metricScheduler(db, &metrics, databaseSizeUpdate, GetDatabeSizeResult, 60*60, 120)
	go metricScheduler(db, &metrics, numberOfConnectionUpdate, GetNumberOfConnectionResult, 5*60, 120)
	log.Printf("Starting Pome %s", Version)
	log.Printf("Application will be available at http://127.0.0.1:%d", *port)
	initWebServer(context)
}
