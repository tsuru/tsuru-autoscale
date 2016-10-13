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

func newAlarm(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var a alarm.Alarm
	err = json.Unmarshal(body, &a)
	if err != nil {
		return err
	}
	err = alarm.NewAlarm(&a)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func listAlarms(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("Authorization")
	alarms, err := alarm.ListAlarmsByToken(token)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(alarms)
}

func listAlarmsByInstance(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	alarms, err := alarm.ListAlarmsByInstance(vars["instance"])
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(alarms)
}

func removeAlarm(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	return alarm.RemoveAlarm(a)
}

func enableAlarm(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	return alarm.Enable(a)
}

func disableAlarm(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	return alarm.Disable(a)
}

func getAlarm(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(a)
}

func listEvents(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	events, err := alarm.EventsByAlarmName(a.Name)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(events)
}
