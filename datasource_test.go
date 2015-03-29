// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"

	"gopkg.in/check.v1"
)

func (s *S) TestRegister(c *check.C) {
	var ds dataSource
	dsFactory := func(conf map[string]interface{}) (dataSource, error) {
		return ds, nil
	}
	Register("graphite", dsFactory)
	d, err := NewDataSource("graphite", nil)
	c.Assert(err, check.IsNil)
	c.Assert(d, check.DeepEquals, ds)
}

type testHandler struct{}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content := `{"name":"Paul"}`
	w.Write([]byte(content))
}

func (s *S) TestHttpDataSourceGet(c *check.C) {
	h := testHandler{}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	ds := httpDataSource{method: "POST", url: ts.URL}
	type dataType struct {
		Name string
	}
	data := dataType{}
	err := ds.Get(&data)
	c.Assert(err, check.IsNil)
	c.Assert(data.Name, check.Equals, "Paul")
}
