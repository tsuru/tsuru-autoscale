// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/http/httptest"

	"github.com/codegangsta/negroni"
	"gopkg.in/check.v1"
)

func (s *S) TestAuthMiddleware(c *check.C) {
	recorder := httptest.NewRecorder()
	a := newAuthMiddleware()
	n := negroni.New()
	n.Use(a)
	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}))
	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	c.Assert(err, check.IsNil)
	n.ServeHTTP(recorder, req)
	c.Assert(recorder.Code, check.Equals, http.StatusUnauthorized)
}

func (s *S) TestAuthMiddlewareIgnoreHealthcheck(c *check.C) {
	recorder := httptest.NewRecorder()
	a := newAuthMiddleware()
	n := negroni.New()
	n.Use(a)
	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}))
	req, err := http.NewRequest("GET", "http://localhost:3000/healthcheck", nil)
	c.Assert(err, check.IsNil)
	n.ServeHTTP(recorder, req)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
}
