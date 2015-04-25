// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tsuru/tsuru-autoscale/tsuru"
)

func serviceAdd(w http.ResponseWriter, r *http.Request) {
	var i tsuru.Instance
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &i)
	if err != nil {
		return
	}
	err = tsuru.NewInstance(&i)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func serviceBindUnit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func serviceBindApp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "{}")
}

func serviceUnbindUnit(w http.ResponseWriter, r *http.Request) {
}

func serviceUnbindApp(w http.ResponseWriter, r *http.Request) {
}

func serviceRemove(w http.ResponseWriter, r *http.Request) {
}
