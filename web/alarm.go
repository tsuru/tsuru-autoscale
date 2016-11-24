// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"fmt"
	"net/http"

	"github.com/ajg/form"
	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/action"
	"github.com/tsuru/tsuru-autoscale/alarm"
	"github.com/tsuru/tsuru-autoscale/datasource"
)

func alarmHandler(w http.ResponseWriter, r *http.Request) error {
	a, err := alarm.FindAlarmBy(nil)
	if err != nil {
		return err
	}
	return render(w, "web/templates/alarm/list.html", a)
}

func alarmDetailHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	ds, err := datasource.FindBy(nil)
	if err != nil {
		return err
	}
	actions, err := action.All()
	if err != nil {
		return err
	}
	context := struct {
		DataSources []datasource.DataSource
		Actions     []action.Action
		Alarm       *alarm.Alarm
	}{
		ds,
		actions,
		a,
	}
	return render(w, "web/templates/alarm/detail.html", context)
}

func alarmAdd(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return err
		}
		ds := []string{}
		for _, d := range r.Form["datasources"] {
			ds = append(ds, d)
		}
		r.Form.Del("datasources")
		envs := map[string]string{}
		for i := range r.Form["key"] {
			if r.Form["key"][i] != "" {
				envs[r.Form["key"][i]] = r.Form["value"][i]
			}
		}
		r.Form.Del("key")
		r.Form.Del("value")
		actions := []string{}
		for _, a := range r.Form["actions"] {
			actions = append(actions, a)
		}
		r.Form.Del("actions")
		var a alarm.Alarm
		d := form.NewDecoder(nil)
		d.IgnoreCase(true)
		d.IgnoreUnknownKeys(true)
		err = d.DecodeValues(&a, r.Form)
		if err != nil {
			return err
		}
		a.DataSources = ds
		a.Actions = actions
		a.Envs = envs
		err = alarm.NewAlarm(&a)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/web/alarm", 302)
		return nil
	}
	ds, err := datasource.FindBy(nil)
	if err != nil {
		return err
	}
	actions, err := action.All()
	if err != nil {
		return err
	}
	context := struct {
		DataSources []datasource.DataSource
		Actions     []action.Action
	}{
		ds,
		actions,
	}
	return render(w, "web/templates/alarm/add.html", context)
}

func alarmRemove(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	err = alarm.RemoveAlarm(a)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/web/alarm", 302)
	return nil
}

func alarmEnable(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	err = alarm.Enable(a)
	if err != nil {
		return err
	}
	http.Redirect(w, r, fmt.Sprintf("/web/alarm/%s", vars["name"]), 302)
	return nil
}

func alarmDisable(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := alarm.FindAlarmByName(vars["name"])
	if err != nil {
		return err
	}
	err = alarm.Disable(a)
	if err != nil {
		return err
	}
	http.Redirect(w, r, fmt.Sprintf("/web/alarm/%s", vars["name"]), 302)
	return nil
}

func alarmEdit(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	ds := []string{}
	for _, d := range r.Form["datasources"] {
		ds = append(ds, d)
	}
	r.Form.Del("datasources")
	envs := map[string]string{}
	for i := range r.Form["key"] {
		if r.Form["key"][i] != "" {
			envs[r.Form["key"][i]] = r.Form["value"][i]
		}
	}
	r.Form.Del("key")
	r.Form.Del("value")
	actions := []string{}
	for _, a := range r.Form["actions"] {
		actions = append(actions, a)
	}
	r.Form.Del("actions")
	d := form.NewDecoder(nil)
	d.IgnoreCase(true)
	d.IgnoreUnknownKeys(true)
	var a alarm.Alarm
	err = d.DecodeValues(&a, r.Form)
	if err != nil {
		return err
	}
	a.DataSources = ds
	a.Actions = actions
	a.Envs = envs
	oldAlarm, err := alarm.FindAlarmByName(a.Name)
	if err != nil {
		return err
	}
	a.Instance = oldAlarm.Instance
	err = alarm.UpdateAlarm(&a)
	if err != nil {
		return err
	}
	u := fmt.Sprintf("/web/alarm/%s", a.Name)
	http.Redirect(w, r, u, 302)
	return nil
}
