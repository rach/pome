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
	web_port = app.Flag("web-port", "web application port (default: 2345)").
			Short('P').Default("2345").PlaceHolder("PORT").Int()
	port = app.Flag("port", "database server port (default: 5432)").
		Short('p').Default("5432").PlaceHolder("PORT").Int()
	sslmode = app.Flag("sslmode", "database SSL mode (default: disable)").
		Short('s').Default("disable").PlaceHolder("SSLMODE").String()
	verbose  = app.Flag("verbose", "").Short('v').Bool()
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
	var connstring = connectionString(*host, *database, *username, pwd, *sslmode, *port)
	db := connectDB(connstring)
	context := &appContext{db, &metrics}
	go metricScheduler(db, &metrics, indexBloatUpdate, GetIndexBloatResult, 12*60*60, 120)
	go metricScheduler(db, &metrics, tableBloatUpdate, GetTableBloatResult, 12*60*60, 120)
	go metricScheduler(db, &metrics, databaseSizeUpdate, GetDatabeSizeResult, 60*60, 120)
	go metricScheduler(db, &metrics, numberOfConnectionUpdate, GetNumberOfConnectionResult, 5*60, 120)
	log.Printf("Starting Pome %s", Version)
	log.Printf("Application will be available at http://127.0.0.1:%d", *web_port)
	initWebServer(context, *web_port)
}
