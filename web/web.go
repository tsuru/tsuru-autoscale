// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func wizard(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/wizard.html")
	t.Execute(w, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func router() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/", index).Methods("GET")
	m.HandleFunc("/wizard", wizard).Methods("GET")
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
