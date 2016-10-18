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
	"github.com/tsuru/tsuru-autoscale/datasource"
	"github.com/tsuru/tsuru-autoscale/wizard"
	"gopkg.in/mgo.v2/bson"
)

func wizardDetailHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/wizard.html")
	if err != nil {
		return err
	}
	vars := mux.Vars(r)
	a, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	return t.Execute(w, a)
}

func wizardHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/wizards.html")
	if err != nil {
		return err
	}
	wizards, err := wizard.FindBy(nil)
	if err != nil {
		return err
	}
	return t.Execute(w, wizards)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/templates/index.html")
	t.Execute(w, nil)
}

func dataSourceHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/datasources.html")
	if err != nil {
		return err
	}
	ds, err := datasource.FindBy(nil)
	if err != nil {
		return err
	}
	return t.Execute(w, ds)
}

func dataSourceDetailHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/datasource.html")
	if err != nil {
		return err
	}
	vars := mux.Vars(r)
	ds, err := datasource.FindBy(bson.M{"name": vars["name"]})
	if err != nil {
		return err
	}
	return t.Execute(w, ds)
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
