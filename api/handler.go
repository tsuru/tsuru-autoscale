// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
)

type authorizationRequiredHandler func(http.ResponseWriter, *http.Request)

func (fn authorizationRequiredHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		err := "Authorization header is required."
		logger().Print(err)
		http.Error(w, err, http.StatusUnauthorized)
	}
	fn(w, r)
}
