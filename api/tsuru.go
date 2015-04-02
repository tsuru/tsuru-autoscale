// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
)

func serviceAdd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func serviceBind(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func serviceUnbind(w http.ResponseWriter, r *http.Request) {
}

func serviceRemove(w http.ResponseWriter, r *http.Request) {
}
