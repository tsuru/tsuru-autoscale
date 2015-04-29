// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/tsuru"
)

func serviceAdd(w http.ResponseWriter, r *http.Request) {
	i := tsuru.Instance{
		Name: r.FormValue("name"),
		Team: r.FormValue("team"),
		User: r.FormValue("user"),
	}
	err := tsuru.NewInstance(&i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func serviceBindUnit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func serviceBindApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := tsuru.GetInstanceByName(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger().Print(err.Error())
		return
	}
	err = i.AddApp(r.FormValue("app-host"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger().Print(err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "{}")
}

func serviceUnbindUnit(w http.ResponseWriter, r *http.Request) {
}

func serviceUnbindApp(w http.ResponseWriter, r *http.Request) {
}

func serviceRemove(w http.ResponseWriter, r *http.Request) {
}
