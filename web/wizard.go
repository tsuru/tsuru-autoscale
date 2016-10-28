// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"fmt"
	"net/http"

	"github.com/ajg/form"
	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/wizard"
)

func wizardDetailHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	return render(w, "web/templates/wizard/detail.html", a)
}

func wizardHandler(w http.ResponseWriter, r *http.Request) error {
	wizards, err := wizard.FindBy(nil)
	if err != nil {
		return err
	}
	return render(w, "web/templates/wizard/list.html", wizards)
}

func wizardAdd(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return err
		}
		var a wizard.AutoScale
		d := form.NewDecoder(nil)
		d.IgnoreCase(true)
		d.IgnoreUnknownKeys(true)
		err = d.DecodeValues(&a, r.Form)
		if err != nil {
			return err
		}
		err = wizard.New(&a)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/web/wizard", 302)
		return nil
	}
	return render(w, "web/templates/wizard/add.html", nil)
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

func wizardEnable(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	err = a.Enable()
	if err != nil {
		return err
	}
	http.Redirect(w, r, fmt.Sprintf("/web/wizard/%s", a.Name), 302)
	return nil
}

func wizardDisable(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	err = a.Disable()
	if err != nil {
		return err
	}
	http.Redirect(w, r, fmt.Sprintf("/web/wizard/%s", a.Name), 302)
	return nil
}
