// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/log"
)

func logger() *log.Logger {
	return log.Log()
}

type handler func(http.ResponseWriter, *http.Request) error

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		logger().Error(err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Router return a http.Handler with all api routes
func Router(m *mux.Router) {
	m.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	m.HandleFunc("/datasource", newDataSource).Methods("POST")
	m.HandleFunc("/datasource", allDataSources).Methods("GET")
	m.HandleFunc("/datasource/{name}", removeDataSource).Methods("DELETE")
	m.HandleFunc("/datasource/{name}", getDataSource).Methods("GET")
	m.HandleFunc("/action", allActions).Methods("GET")
	m.HandleFunc("/action", newAction).Methods("POST")
	m.HandleFunc("/action/{name}", removeAction).Methods("DELETE")
	m.HandleFunc("/action/{name}", actionInfo).Methods("GET")
	m.Handle("/alarm", handler(newAlarm)).Methods("POST")
	m.Handle("/alarm/instance/{instance}", handler(listAlarmsByInstance)).Methods("GET")
	m.Handle("/alarm", authorizationRequiredHandler(listAlarms)).Methods("GET")
	m.Handle("/alarm/{name}/enable", handler(enableAlarm)).Methods("PUT")
	m.Handle("/alarm/{name}/disable", handler(disableAlarm)).Methods("PUT")
	m.Handle("/alarm/{name}", handler(removeAlarm)).Methods("DELETE")
	m.Handle("/alarm/{name}", handler(getAlarm)).Methods("GET")
	m.Handle("/alarm/{name}/event", handler(listEvents)).Methods("GET")
	m.HandleFunc("/resources", serviceAdd)
	m.HandleFunc("/resources/{name}/bind", serviceBindUnit).Methods("POST")
	m.HandleFunc("/resources/{name}/bind-app", serviceBindApp).Methods("POST")
	m.HandleFunc("/resources/{name}/bind-app", serviceUnbindApp).Methods("DELETE")
	m.HandleFunc("/resources/{name}/bind", serviceUnbindUnit).Methods("DELETE")
	m.HandleFunc("/resources/{name}", serviceRemove).Methods("DELETE")
	m.HandleFunc("/service/instance/{name}", serviceInstanceByName).Methods("GET")
	m.Handle("/service/instance", authorizationRequiredHandler(serviceInstances)).Methods("GET")
	m.HandleFunc("/wizard/{name}/events", eventsByWizardName).Methods("GET")
	m.HandleFunc("/wizard/{name}", wizardByName).Methods("GET")
	m.HandleFunc("/wizard/{name}", removeWizard).Methods("DELETE")
	m.HandleFunc("/wizard", newAutoScale).Methods("POST")
}
