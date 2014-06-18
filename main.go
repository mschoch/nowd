package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":4793", "bind address")
var timeout = flag.Duration("timeout", 5*time.Minute, "expire sensor readings after this duration")
var expirationInterval = flag.Duration("expire", 1*time.Minute, "check expiration interval")
var sensorCache *TimeRevCache

func main() {

	flag.Parse()

	sensorCache = NewTimeRevCache(*timeout)
	go ExpireCache()

	r := mux.NewRouter()
	r.HandleFunc("/", serveRoot)
	r.HandleFunc("/{device}", serveSensorUpdate).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	mustEncode(w, sensorCache.Values())
}

func serveSensorUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	rev := 0
	if r.FormValue("rev") != "" {
		var err error
		rev, err = strconv.Atoi(r.FormValue("rev"))
		if err != nil {
			showError(w, r, fmt.Sprintf("error: %v", err), 400)
			return
		}
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		showError(w, r, fmt.Sprintf("error: %v", err), 400)
		return
	}

	var message interface{}
	err = json.Unmarshal(bodyBytes, &message)
	if err != nil {
		showError(w, r, fmt.Sprintf("error parsing JSON: %v bytes: `%s`", err, bodyBytes), 400)
		return
	}

	updated := sensorCache.CheckAndUpdate(device, rev, message)
	if updated {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNotModified)
	}
}

func ExpireCache() {
	ticker := time.NewTicker(*expirationInterval)
	for _ = range ticker.C {
		sensorCache.Expire()
	}
}
