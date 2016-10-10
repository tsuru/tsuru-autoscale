// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// Router return a http.Handler with all web routes
func Router() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/", indexHandler).Methods("GET")
	m.HandleFunc("/wizard", wizardHandler).Methods("GET")
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(m)
	return n
}
