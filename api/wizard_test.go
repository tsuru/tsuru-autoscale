// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"gopkg.in/check.v1"
)

func (s *S) TestNewAutoScale(c *check.C) {
	body := `{"name":"test","minUnits":2,"scaleUp":{},"scaleDown":{}}`
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/wizard", strings.NewReader(body))
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}
