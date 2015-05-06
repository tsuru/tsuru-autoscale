// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	stdlog "log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/log"
)

func logger() *stdlog.Logger {
	return log.Logger()
}

func Router() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	m.HandleFunc("/datasource", newDataSource).Methods("POST")
	m.HandleFunc("/datasource", allDataSources).Methods("GET")
	m.HandleFunc("/action", allActions).Methods("GET")
	m.HandleFunc("/action", newAction).Methods("POST")
	m.HandleFunc("/alarm", newAlarm).Methods("POST")
	m.HandleFunc("/resources", serviceAdd)
	m.HandleFunc("/resources/{name}/bind", serviceBindUnit).Methods("POST")
	m.HandleFunc("/resources/{name}/bind-app", serviceBindApp).Methods("POST")
	m.HandleFunc("/resources/{name}/bind-app", serviceUnbindApp).Methods("DELETE")
	m.HandleFunc("/resources/{name}/bind", serviceUnbindUnit).Methods("DELETE")
	m.HandleFunc("/resources/{name}", serviceRemove).Methods("DELETE")
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(m)
	return n
}
