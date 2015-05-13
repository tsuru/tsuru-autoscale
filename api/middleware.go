// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import "net/http"

// authMiddleware is a middleware handler that checks if the Authorization header exists.
type authMiddleware struct{}

// newAuthMiddleware returns a new AuthMiddleware instance.
func newAuthMiddleware() *authMiddleware {
	return &authMiddleware{}
}

func (a *authMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path != "/healthcheck" {
		token := r.Header.Get("Authorization")
		if token == "" {
			err := "Authorization header is required."
			http.Error(rw, err, http.StatusUnauthorized)
			logger().Print(err)
		}
	}
	next(rw, r)
}
