// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/wizard"
)

func wizardDetailHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/wizard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	a, err := wizard.FindByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = t.Execute(w, a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func wizardHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/wizards.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wizards, err := wizard.FindBy(nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, wizards)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/templates/index.html")
	t.Execute(w, nil)
}
