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

func newAutoScale(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	var a wizard.AutoScale
	err = json.Unmarshal(body, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	err = wizard.New(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}

func wizardByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	autoScale, err := wizard.FindByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&autoScale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}

func removeWizard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	autoScale, err := wizard.FindByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
	}
	err = wizard.Remove(autoScale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}

func eventsByWizardName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	autoScale, err := wizard.FindByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	events, err := autoScale.Events()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
	}
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}
