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
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	var a alarm.Alarm
	err = json.Unmarshal(body, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	err = alarm.NewAlarm(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}

func listAlarms(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	alarms, err := alarm.ListAlarmsByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(alarms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}

func listAlarmsByInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alarms, err := alarm.ListAlarmsByInstance(vars["instance"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(alarms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}

func removeAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
	}
	err = alarm.RemoveAlarm(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}

func enableAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
	}
	err = alarm.Enable(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}

func disableAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
	}
	err = alarm.Disable(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}

func getAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
	}
	events, err := alarm.EventsByAlarmName(a.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}
