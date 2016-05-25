package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rach/pome/Godeps/_workspace/src/github.com/elazarl/go-bindata-assetfs"
	"github.com/rach/pome/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	_ "github.com/rach/pome/Godeps/_workspace/src/github.com/lib/pq"
)

type appContext struct {
	db      *sqlx.DB
	metrics *MetricList
}

type appHandler struct {
	*appContext
	H func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func metricsHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	js, err := json.Marshal(a.metrics)
	if err != nil {
		return 500, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return 200, nil
}

func aliasHandler(rw http.ResponseWriter, req *http.Request) {
	if bs, err := Asset("index.html"); err != nil {
		rw.WriteHeader(http.StatusNotFound)
	} else {
		var reader = bytes.NewBuffer(bs)
		io.Copy(rw, reader)
	}
}

func newStaticHandler() http.Handler {
	fs := assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}
	return http.FileServer(&fs)
}

func initWebServer(context *appContext, webHost string, webPort int) {
	http.Handle("/api/stats", appHandler{context, metricsHandler})
	http.HandleFunc("/about", aliasHandler)
	http.HandleFunc("/bloat/indexes", aliasHandler)
	http.HandleFunc("/bloat/tables", aliasHandler)
	http.Handle("/", newStaticHandler())
	http.ListenAndServe(fmt.Sprintf("%s:%d", webHost, webPort), nil)
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
