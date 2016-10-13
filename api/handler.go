// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
)

type authorizationRequiredHandler func(http.ResponseWriter, *http.Request) error

func (fn authorizationRequiredHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		msg := "Authorization header is required."
		logger().Print(msg)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	err := fn(w, r)
	if err != nil {
		logger().Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
