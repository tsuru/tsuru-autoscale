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
