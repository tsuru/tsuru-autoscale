// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/alarm"
)

func newAlarm(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var a alarm.Alarm
	err = json.Unmarshal(body, &a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = alarm.NewAlarm(&a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
}

func listAlarms(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	alarms, err := alarm.ListAlarmsByToken(token)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(alarms)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func listAlarmsByInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alarms, err := alarm.ListAlarmsByInstance(vars["instance"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(alarms)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func removeAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	err = alarm.RemoveAlarm(a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func enableAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	err = alarm.Enable(a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func disableAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	err = alarm.Disable(a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	events, err := alarm.EventsByAlarmName(a.Name)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
