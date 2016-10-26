// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/wizard"
)

func newAutoScale(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var a wizard.AutoScale
	err = json.Unmarshal(body, &a)
	if err != nil {
		return err
	}
	err = wizard.New(&a)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func wizardByName(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	autoScale, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(&autoScale)
}

func removeWizard(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	autoScale, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	return wizard.Remove(autoScale)
}

func eventsByWizardName(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	autoScale, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	events, err := autoScale.Events()
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(&events)
}

func wizardEnable(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	autoScale, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	return autoScale.Enable()
}

func wizardDisable(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	autoScale, err := wizard.FindByName(vars["name"])
	if err != nil {
		return err
	}
	return autoScale.Disable()
}

func wizardUpdate(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var a wizard.AutoScale
	err = json.Unmarshal(body, &a)
	if err != nil {
		return err
	}
	vars := mux.Vars(r)
	a.Name = vars["name"]
	err = wizard.Update(&a)
	if err != nil {
		return err
	}
	return nil
}
