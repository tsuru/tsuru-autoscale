// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/action"
	"github.com/tsuru/tsuru-autoscale/alarm"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/templates/index.html")
	t.Execute(w, nil)
}

func actionHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/actions.html")
	if err != nil {
		return err
	}
	a, err := action.All()
	if err != nil {
		return err
	}
	return t.Execute(w, a)
}

func actionDetailHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/action.html")
	if err != nil {
		return err
	}
	vars := mux.Vars(r)
	a, err := action.FindByName(vars["name"])
	if err != nil {
		return err
	}
	return t.Execute(w, a)
}

func alarmHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/alarms.html")
	if err != nil {
		return err
	}
	a, err := alarm.FindAlarmBy(nil)
	if err != nil {
		return err
	}
	return t.Execute(w, a)
}

func alarmDetailHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/alarm.html")
	if err != nil {
		return err
	}
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	return t.Execute(w, a)
}
