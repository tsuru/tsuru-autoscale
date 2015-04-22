// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"gopkg.in/check.v1"
)

func (s *S) TestNewAction(c *check.C) {
	body := `{"name":"new","url":"http://tsuru.io","method":"GET"}`
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/action", strings.NewReader(body))
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	fmt.Println(recorder)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}
