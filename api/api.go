// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/log"
)

func logger() *log.Logger {
	return log.Log()
}

// Router return a http.Handler with all api routes 
func Router() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	m.HandleFunc("/datasource", newDataSource).Methods("POST")
	m.HandleFunc("/datasource", allDataSources).Methods("GET")
	m.HandleFunc("/datasource/{name}", removeDataSource).Methods("DELETE")
	m.HandleFunc("/datasource/{name}", getDataSource).Methods("GET")
	m.HandleFunc("/action", allActions).Methods("GET")
	m.HandleFunc("/action", newAction).Methods("POST")
	m.HandleFunc("/action/{name}", removeAction).Methods("DELETE")
	m.HandleFunc("/action/{name}", actionInfo).Methods("GET")
	m.HandleFunc("/alarm", newAlarm).Methods("POST")
	m.HandleFunc("/alarm/instance/{instance}", listAlarmsByInstance).Methods("GET")
	m.HandleFunc("/alarm", listAlarms).Methods("GET")
	m.HandleFunc("/alarm/{name}/enable", enableAlarm).Methods("PUT")
	m.HandleFunc("/alarm/{name}/disable", disableAlarm).Methods("PUT")
	m.HandleFunc("/alarm/{name}", removeAlarm).Methods("DELETE")
	m.HandleFunc("/alarm/{name}", getAlarm).Methods("GET")
	m.HandleFunc("/alarm/{name}/event", listEvents).Methods("GET")
	m.HandleFunc("/resources", serviceAdd)
	m.HandleFunc("/resources/{name}/bind", serviceBindUnit).Methods("POST")
	m.HandleFunc("/resources/{name}/bind-app", serviceBindApp).Methods("POST")
	m.HandleFunc("/resources/{name}/bind-app", serviceUnbindApp).Methods("DELETE")
	m.HandleFunc("/resources/{name}/bind", serviceUnbindUnit).Methods("DELETE")
	m.HandleFunc("/resources/{name}", serviceRemove).Methods("DELETE")
	m.HandleFunc("/service/instance/{name}", serviceInstanceByName).Methods("GET")
	m.HandleFunc("/service/instance", serviceInstances).Methods("GET")
	m.HandleFunc("/wizard/{name}/events", eventsByWizardName).Methods("GET")
	m.HandleFunc("/wizard/{name}", wizardByName).Methods("GET")
	m.HandleFunc("/wizard/{name}", removeWizard).Methods("DELETE")
	m.HandleFunc("/wizard", newAutoScale).Methods("POST")
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.Use(newAuthMiddleware())
	n.UseHandler(m)
	return n
}
