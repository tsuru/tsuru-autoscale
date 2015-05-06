// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tsuru/tsuru-autoscale/datasource"
)

func newDataSource(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	var ds datasource.DataSource
	err = json.Unmarshal(body, &ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	err = datasource.New(&ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}

func allDataSources(w http.ResponseWriter, r *http.Request) {
	ds, err := datasource.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
	}
}
