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

func newAction(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var a action.Action
	err = json.Unmarshal(body, &a)
	if err != nil {
		return err
	}
	err = action.New(&a)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func allActions(w http.ResponseWriter, r *http.Request) error {
	actions, err := action.All()
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(actions)
}

func removeAction(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := action.FindByName(vars["name"])
	if err != nil {
		return err
	}
	return action.Remove(a)
}

func actionInfo(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	a, err := action.FindByName(vars["name"])
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(a)
}
