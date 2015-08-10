// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tsuru/tsuru-autoscale/wizard"
)

func newAutoScale(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	var a wizard.AutoScale
	err = json.Unmarshal(body, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	err = wizard.New(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger().Print(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}
