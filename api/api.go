// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func Router() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/datasource", newDataSource).Methods("POST")
	m.HandleFunc("/action", newAction).Methods("POST")
	m.HandleFunc("/resources", serviceAdd)
	m.HandleFunc("/resources/{name}/bind", serviceBind).Methods("POST")
	m.HandleFunc("/resources/{name}/bind", serviceUnbind).Methods("DELETE")
	m.HandleFunc("/resources/{name}", serviceUnbind).Methods("DELETE")
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(m)
	return n
}
