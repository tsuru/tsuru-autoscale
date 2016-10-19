// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/datasource"
	"gopkg.in/mgo.v2/bson"
)

func dataSourceHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/datasources.html")
	if err != nil {
		return err
	}
	ds, err := datasource.FindBy(nil)
	if err != nil {
		return err
	}
	return t.Execute(w, ds)
}

func dataSourceDetailHandler(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("web/templates/datasource.html")
	if err != nil {
		return err
	}
	vars := mux.Vars(r)
	ds, err := datasource.FindBy(bson.M{"name": vars["name"]})
	if err != nil {
		return err
	}
	return t.Execute(w, ds)
}
