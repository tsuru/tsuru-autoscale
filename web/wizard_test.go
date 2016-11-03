// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"

	"github.com/tsuru/tsuru-autoscale/wizard"
	"gopkg.in/check.v1"
)

func (s *S) TestWizardRemove(c *check.C) {
	recorder := httptest.NewRecorder()
	a := wizard.AutoScale{Name: "new"}
	err := wizard.New(&a)
	c.Assert(err, check.IsNil)
	request, err := http.NewRequest("GET", "/wizard/new/delete", nil)
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
	_, err = wizard.FindByName(a.Name)
	c.Assert(err, check.NotNil)
}
