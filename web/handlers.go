// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"html/template"
	"net/http"
)

func wizardHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/templates/wizard.html")
	t.Execute(w, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/templates/index.html")
	t.Execute(w, nil)
}
