// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"

	"github.com/ajg/form"
	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/datasource"
	"gopkg.in/mgo.v2/bson"
)

func dataSourceHandler(w http.ResponseWriter, r *http.Request) error {
	ds, err := datasource.FindBy(nil)
	if err != nil {
		return err
	}
	return render(w, "web/templates/datasource/list.html", ds)
}

func dataSourceDetailHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ds, err := datasource.FindBy(bson.M{"name": vars["name"]})
	if err != nil {
		return err
	}
	return render(w, "web/templates/datasource/detail.html", ds[0])
}

func dataSourceAdd(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
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
		headers := map[string]string{}
		for i := range r.Form["key"] {
			if r.Form["key"][i] != "" {
				headers[r.Form["key"][i]] = r.Form["value"][i]
			}
		}
		ds.Headers = headers
		err = datasource.New(&ds)
		if err != nil {
			return err
		}
		http.Redirect(w, r, "/web/datasource", 302)
		return nil
	}
	return render(w, "web/templates/datasource/add.html", nil)
}

func dataSourceRemoveHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ds, err := datasource.FindBy(bson.M{"name": vars["name"]})
	if err != nil {
		return err
	}
	err = datasource.Remove(&ds[0])
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/web/datasource", 302)
	return nil
}
