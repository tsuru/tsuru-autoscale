// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type handler func(http.ResponseWriter, *http.Request) error

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Router return a http.Handler with all web routes
func Router(m *mux.Router) {
	m.HandleFunc("/", indexHandler).Methods("GET")
	m.Handle("/event", handler(eventHandler)).Methods("GET")
	m.Handle("/alarm", handler(alarmHandler)).Methods("GET")
	m.Handle("/alarm/add", handler(alarmAdd)).Methods("GET", "POST")
	m.Handle("/alarm/{name}", handler(alarmDetailHandler)).Methods("GET")
	m.Handle("/alarm/{name}/delete", handler(alarmRemove)).Methods("GET")
	m.Handle("/alarm/{name}/enable", handler(alarmEnable)).Methods("GET")
	m.Handle("/alarm/{name}/disable", handler(alarmDisable)).Methods("GET")
	m.Handle("/alarm/{name}/edit", handler(alarmEdit)).Methods("POST")
	m.Handle("/action", handler(actionHandler)).Methods("GET")
	m.Handle("/action/add", handler(actionAdd)).Methods("GET", "POST")
	m.Handle("/action/{name}", handler(actionDetailHandler)).Methods("GET")
	m.Handle("/action/{name}/delete", handler(actionRemove)).Methods("GET")
	m.Handle("/datasource", handler(dataSourceHandler)).Methods("GET")
	m.Handle("/datasource/add", handler(dataSourceAdd)).Methods("GET", "POST")
	m.Handle("/datasource/{name}", handler(dataSourceDetailHandler)).Methods("GET")
	m.Handle("/datasource/{name}/delete", handler(dataSourceRemoveHandler)).Methods("GET")
	m.Handle("/wizard", handler(wizardHandler)).Methods("GET")
	m.Handle("/wizard/{name}", handler(wizardDetailHandler)).Methods("GET")
	m.Handle("/wizard/{name}/delete", handler(wizardRemove)).Methods("GET")
	m.Handle("/wizard/{name}/enable", handler(wizardEnable)).Methods("GET")
	m.Handle("/wizard/{name}/disable", handler(wizardDisable)).Methods("GET")
}
