// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"github.com/gorilla/mux"
)

// Router return a http.Handler with all web routes
func Router(m *mux.Router) {
	m.HandleFunc("/", indexHandler).Methods("GET")
	m.HandleFunc("/datasource", dataSourceHandler).Methods("GET")
	m.HandleFunc("/wizard", wizardHandler).Methods("GET")
	m.HandleFunc("/wizard/{name}", wizardDetailHandler).Methods("GET")
}
