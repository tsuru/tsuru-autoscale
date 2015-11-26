// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/tsuru"
	"github.com/tsuru/tsuru-autoscale/wizard"
)

func serviceAdd(w http.ResponseWriter, r *http.Request) {
	i := tsuru.Instance{
		Name: r.FormValue("name"),
		Team: r.FormValue("team"),
		User: r.FormValue("user"),
	}
	err := tsuru.NewInstance(&i)
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = i.AddApp(r.FormValue("app-host"))
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "{}")
}

func serviceUnbindUnit(w http.ResponseWriter, r *http.Request) {
}

func serviceUnbindApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := tsuru.GetInstanceByName(vars["name"])
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	r.Method = "POST"
	err = i.RemoveApp(r.FormValue("app-host"))
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	autoScale, err := wizard.FindByName(vars["name"])
	if err == nil {
		rerr := wizard.Remove(autoScale)
		if rerr != nil {
			logger().Error(rerr.Error())
			http.Error(w, rerr.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func serviceRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := tsuru.GetInstanceByName(vars["name"])
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = tsuru.RemoveInstance(i)
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func serviceInstances(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	instances, err := tsuru.FindServiceInstance(token)
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&instances)
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serviceInstanceByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance, err := tsuru.GetInstanceByName(vars["name"])
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&instance)
	if err != nil {
		logger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
