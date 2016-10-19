// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"net/http"

	"github.com/ajg/form"
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

func dataSourceAdd(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			return err
		}
		var ds datasource.DataSource
		d := form.NewDecoder(nil)
		d.IgnoreCase(true)
		d.IgnoreUnknownKeys(true)
		err = d.DecodeValues(&ds, r.Form)
		if err != nil {
			return err
		}
		err = datasource.New(&ds)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/web/datasource", 302)
		return nil
	}
	t, err := template.ParseFiles("web/templates/datasource-add.html")
	if err != nil {
		return err
	}
	return t.Execute(w, nil)
}
