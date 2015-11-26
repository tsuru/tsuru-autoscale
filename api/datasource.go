// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/datasource"
	"gopkg.in/mgo.v2/bson"
)

func newDataSource(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var ds datasource.DataSource
	err = json.Unmarshal(body, &ds)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = datasource.New(&ds)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
}

func allDataSources(w http.ResponseWriter, r *http.Request) {
	var q bson.M
	public, err := strconv.ParseBool(r.URL.Query().Get("public"))
	if err == nil {
		q = bson.M{"public": public}
	}
	ds, err := datasource.FindBy(q)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ds)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func removeDataSource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds, err := datasource.Get(vars["name"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	err = datasource.Remove(ds)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getDataSource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ds, err := datasource.Get(vars["name"])
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ds)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
