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

func newDataSource(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var ds datasource.DataSource
	err = json.Unmarshal(body, &ds)
	if err != nil {
		return err
	}
	err = datasource.New(&ds)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func allDataSources(w http.ResponseWriter, r *http.Request) error {
	var q bson.M
	public, err := strconv.ParseBool(r.URL.Query().Get("public"))
	if err == nil {
		q = bson.M{"public": public}
	}
	ds, err := datasource.FindBy(q)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(ds)
}

func removeDataSource(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ds, err := datasource.Get(vars["name"])
	if err != nil {
		return err
	}
	return datasource.Remove(ds)
}

func getDataSource(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ds, err := datasource.Get(vars["name"])
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(ds)
}
