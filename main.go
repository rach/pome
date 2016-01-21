package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/rach/pome/Godeps/_workspace/src/github.com/alecthomas/kingpin"
	"github.com/robfig/cron"
)

//go:generate go-bindata -prefix "static/" -pkg main -o bindata.go static/index.html static/favicons/... static/build/...

const (
	Version = "0.2.0"
)

type CronValue string

func (c *CronValue) Set(value string) error {
	_, err := cron.Parse(value)
	if err != nil {
		return fmt.Errorf("expected cron expression or '@every DURATION' got '%s'", value)
	}
	*c = (CronValue)(value)
	return nil
}

func (c *CronValue) String() string {
	return (string)(*c)
}

func CronFlag(s kingpin.Settings) (target *string) {
	target = new(string)
	s.SetValue((*CronValue)(target))
	return
}

func addUsernameFlag(app *kingpin.Application) *string {
	uname := os.Getenv("PGUSER")
	if uname == "" {
		u, err := user.Current()
		if err != nil {
			return app.Flag("username", "").Short('U').
				PlaceHolder("USERNAME").Required().String()
		}
		uname = u.Username
	}
	return app.Flag("username", "").Short('U').Default(uname).
		PlaceHolder(fmt.Sprintf("USERNAME (default: %s)", uname)).String()
}

var (
	app  = kingpin.New("pome", "A Postgres Metrics Dashboard.")
	host = app.Flag("host", "database server host (default: localhost)").
		OverrideDefaultFromEnvar("PGHOST").
		Short('h').PlaceHolder("HOSTNAME").String()
	webHost = app.Flag("web-host", "web application host (default: 127.0.0.1)").
		Short('H').Default("127.0.0.1").PlaceHolder("WEBHOST").String()
	webPort = app.Flag("web-port", "web application port (default: 2345)").
		Short('P').Default("2345").PlaceHolder("WEBPORT").Int()
	port = app.Flag("port", "database server port (default: 5432)").
		Short('p').Default("5432").
		OverrideDefaultFromEnvar("PGPORT").
		PlaceHolder("PORT").Int()
	sslmode = app.Flag("sslmode", "database SSL mode (default: disable)").
		Short('s').Default("disable").PlaceHolder("SSLMODE").String()
	verbose  = app.Flag("verbose", "").Short('v').Bool()
	password = app.Flag("password", "").Short('W').Bool()
	username = addUsernameFlag(app)
	database = app.Arg("DBNAME", "").Required().String()
	// Scheduling flags
	scheduleTableBloat = CronFlag(
		app.
			Flag("schedule-table-bloat", "Cron like expression for when to query table bloat").
			Default("@every 12h"),
	)
	scheduleIndexBloat = CronFlag(
		app.
			Flag("schedule-index-bloat", "Cron like expression for when to the query index bloat").
			Default("@every 12h"),
	)
	scheduleDbSize = CronFlag(
		app.
			Flag("schedule-db-size", "Cron like expression for when to query the database size").
			Default("@every 1h"),
	)
	scheduleNumConn = CronFlag(
		app.
			Flag("schedule-num-conn", "Cron like expression for when to query the number of connection").
			Default("@every 5m"),
	)
	scheduleTPS = CronFlag(
		app.
			Flag("schedule-tps", "Cron like expression for when to query the transaction per second estimate").
			Default("@every 1m"),
	)
)

func parseCmdLine(args []string) (command string, err error) {
	//this is isolated from the main() function to make it more testable
	app.Version(Version)
	return app.Parse(args)
}

func main() {
	kingpin.MustParse(parseCmdLine(os.Args[1:]))
	var metrics = initMetricList(Version)
	pwd := os.Getenv("PGPASSWORD")
	if *password {
		fmt.Print("Enter Password: ")
		fmt.Scanln(&pwd)
	}
	db := connectDB(*host, *database, *username, pwd, *sslmode, *port)
	context := &appContext{db, &metrics}
	log.Printf("Starting Pome %s", Version)
	log.Printf("Application will be available at http://%s:%d", *webHost, *webPort)

	initScheduler(
		db,
		&metrics,
		*scheduleTableBloat,
		*scheduleDbSize,
		*scheduleIndexBloat,
		*scheduleNumConn,
		*scheduleTPS,
	)
	initWebServer(context, *webHost, *webPort)
}
