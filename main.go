// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/alarm"
	"github.com/tsuru/tsuru-autoscale/api"
	"github.com/tsuru/tsuru-autoscale/web"
)

func port() string {
	var p string
	if p = os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func router() http.Handler {
	m := mux.NewRouter()
	apiRouter := m.PathPrefix("/").Subrouter()
	api.Router(apiRouter)
	webRouter := m.PathPrefix("/web").Subrouter()
	web.Router(webRouter)
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(m)
	return n
}

func runServer() {
	http.Handle("/", router())
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port()), nil))
}

func main() {
	if os.Args[1] == "agent" {
		alarm.StartAutoScale()
	} else {
		runServer()
	}
}
