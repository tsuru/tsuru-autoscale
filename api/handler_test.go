// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/http/httptest"

	"github.com/codegangsta/negroni"
	"gopkg.in/check.v1"
)

func (s *S) TestAuthHandlerWithoutToken(c *check.C) {
	recorder := httptest.NewRecorder()
	n := negroni.New()
	n.UseHandler(authorizationRequiredHandler(func(rw http.ResponseWriter, r *http.Request) {}))
	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	c.Assert(err, check.IsNil)
	n.ServeHTTP(recorder, req)
	c.Assert(recorder.Code, check.Equals, http.StatusUnauthorized)
}

func (s *S) TestAuthHandler(c *check.C) {
	recorder := httptest.NewRecorder()
	n := negroni.New()
	n.UseHandler(authorizationRequiredHandler(func(rw http.ResponseWriter, r *http.Request) {}))
	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	req.Header.Add("Authorization", "1234")
	c.Assert(err, check.IsNil)
	n.ServeHTTP(recorder, req)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
}
