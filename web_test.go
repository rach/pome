package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler(t *testing.T) {
	m := MetricList{Version: "1"}
	a := appContext{nil, &m}
	w := httptest.NewRecorder()
	r := http.Request{}

	status, _ := metricsHandler(&a, w, &r)

	if status != 200 {
		t.Error("Response should be successful")
	}

	if w.HeaderMap["Content-Type"][0] != "application/json" {
		t.Error("Content-Type should be set to application/json")
	}

	v := MetricList{}
	json.Unmarshal(w.Body.Bytes(), &v)

	if v.Version != m.Version {
		t.Errorf("Expected %v to be %v", v, m)
	}
}

func TestAliasHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := http.Request{}

	aliasHandler(w, &r)
	bs, _ := Asset("index.html")
	html := string(bs)

	if w.Body.String() != html {
		t.Error("Index page should be rendered")
	}
}
