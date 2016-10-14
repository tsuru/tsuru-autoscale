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

func serviceAdd(w http.ResponseWriter, r *http.Request) error {
	i := tsuru.Instance{
		Name: r.FormValue("name"),
		Team: r.FormValue("team"),
		User: r.FormValue("user"),
	}
	err := tsuru.NewInstance(&i)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func serviceBindUnit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func serviceBindApp(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	i, err := tsuru.GetInstanceByName(vars["name"])
	if err != nil {
		return err
	}
	err = i.AddApp(r.FormValue("app-host"))
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "{}")
	return nil
}

func serviceUnbindUnit(w http.ResponseWriter, r *http.Request) {
}

func serviceUnbindApp(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	i, err := tsuru.GetInstanceByName(vars["name"])
	if err != nil {
		return err
	}
	r.Method = "POST"
	err = i.RemoveApp(r.FormValue("app-host"))
	if err != nil {
		return err
	}
	autoScale, err := wizard.FindByName(vars["name"])
	if err == nil {
		rerr := wizard.Remove(autoScale)
		if rerr != nil {
			return rerr
		}
	}
	return nil
}

func serviceRemove(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	i, err := tsuru.GetInstanceByName(vars["name"])
	if err != nil {
		return nil
	}
	return tsuru.RemoveInstance(i)
}

func serviceInstances(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("Authorization")
	instances, err := tsuru.FindServiceInstance(token)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(&instances)
}

func serviceInstanceByName(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	instance, err := tsuru.GetInstanceByName(vars["name"])
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(&instance)
}
