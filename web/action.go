// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"net/http"

	"github.com/ajg/form"
	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/action"
)

func actionHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/action/list.html")
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
	t, err := template.ParseFiles("web/templates/action/detail.html")
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

func actionAdd(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return err
		}
		var a action.Action
		d := form.NewDecoder(nil)
		d.IgnoreCase(true)
		d.IgnoreUnknownKeys(true)
		err = d.DecodeValues(&a, r.Form)
		if err != nil {
			return err
		}
		err = action.New(&a)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/web/action", 302)
		return nil
	}
	t, err := template.ParseFiles("web/templates/action/add.html")
	if err != nil {
		return err
	}
	return t.Execute(w, nil)
}

func actionRemove(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := action.FindByName(vars["name"])
	if err != nil {
		return err
	}
	err = action.Remove(a)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/web/action", 302)
	return nil
}
