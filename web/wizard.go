// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"net/http"

	"github.com/ajg/form"
	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/wizard"
)

func wizardDetailHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/wizard/detail.html")
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
	t, err := template.ParseFiles("web/templates/wizard/list.html")
	if err != nil {
		return err
	}
	wizards, err := wizard.FindBy(nil)
	if err != nil {
		return err
	}
	return t.Execute(w, wizards)
}

func wizardAdd(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return err
		}
		var w wizard.AutoScale
		d := form.NewDecoder(nil)
		d.IgnoreCase(true)
		d.IgnoreUnknownKeys(true)
		err = d.DecodeValues(&w, r.Form)
		if err != nil {
			return err
		}
		err = wizard.New(&w)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/web/wizard", 302)
		return nil
	}
	t, err := template.ParseFiles("web/templates/wizard/add.html")
	if err != nil {
		return err
	}
	return t.Execute(w, nil)
}

func wizardRemove(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	err = wizard.Remove(a)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/web/wizard", 302)
	return nil
}
