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
)

func router() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/", indexHandler).Methods("GET")
	m.HandleFunc("/wizard", wizardHandler).Methods("GET")
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(m)
	return n
}

func port() string {
	var p string
	if p = os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func runServer() {
	http.Handle("/", router())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port()), nil))
}

func main() {
	runServer()
}
