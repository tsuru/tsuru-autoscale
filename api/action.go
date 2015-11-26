// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/action"
)

func newAction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var a action.Action
	err = json.Unmarshal(body, &a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = action.New(&a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
}

func allActions(w http.ResponseWriter, r *http.Request) {
	actions, err := action.All()
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(actions)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func removeAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := action.FindByName(vars["name"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	err = action.Remove(a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func actionInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := action.FindByName(vars["name"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(a)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
