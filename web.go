package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io"
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

func staticHandler(rw http.ResponseWriter, req *http.Request) {
	var path string = req.URL.Path
	//fmt.Println(path)
	if path == "" {
		path = "index.html"
	}
	if bs, err := Asset(path); err != nil {
		rw.WriteHeader(http.StatusNotFound)
	} else {
		fmt.Fprintf(rw, http.DetectContentType(bs[:]))
		rw.Header().Set("Content-Type", http.DetectContentType(bs[:]))
		var reader = bytes.NewBuffer(bs)
		io.Copy(rw, reader)
	}
}

type appContext struct {
	db      *sqlx.DB
	metrics *MetricList
}

type appHandler struct {
	*appContext
	H func(*appContext, http.ResponseWriter, *http.Request) (int, error)
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
