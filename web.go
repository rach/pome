package main

import (
	"encoding/json"
	"fmt"
	"github.com/rach/poda/Godeps/_workspace/src/github.com/elazarl/go-bindata-assetfs"
	"github.com/rach/poda/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	_ "github.com/rach/poda/Godeps/_workspace/src/github.com/lib/pq"
	"log"
	"net/http"
)

func metricsHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	// Our handlers now have access to the members of our context struct.
	// e.g. we can call methods on our DB type via err := a.db.GetPosts()
	v, _ := json.Marshal(a.metrics)
	fmt.Fprintf(w, string(v))
	return 200, nil
}

type appContext struct {
	db      *sqlx.DB
	metrics *MetricList
}

type appHandler struct {
	*appContext
	H func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func initWebServer(context *appContext) {
	http.Handle("/api/stats", appHandler{context, metricsHandler})
	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.appContext as a parameter to our handler type.
	status, err := ah.H(ah.appContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page, we can
			// now leverage our context instance - e.g.
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}
