// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"net/http"
	"net/http/httptest"
	"os"

	"gopkg.in/check.v1"
)

func (s *S) TestFindServiceInstance(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"Name":"instance"}]`))
	}))
	defer ts.Close()
	err := os.Setenv("TSURU_HOST", ts.URL)
	c.Assert(err, check.IsNil)
	token := "token"
	instances, err := FindServiceInstance(token)
	c.Assert(err, check.IsNil)
	c.Assert(instances, check.HasLen, 1)
	c.Assert(instances[0].Name, check.Equals, "instance")
}
