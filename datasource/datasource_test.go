// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package datasource

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	conn *db.Storage
}

func (s *S) SetUpSuite(c *check.C) {
	var err error
	s.conn, err = db.Conn()
	c.Assert(err, check.IsNil)
}

var _ = check.Suite(&S{})

type testHandler struct{}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content := `{"name":"Paul"}`
	w.Write([]byte(content))
}

func (s *S) TestHttpDataSourceGet(c *check.C) {
	h := testHandler{}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	ds := DataSource{Method: "POST", URL: ts.URL}
	type dataType struct {
		Name string
	}
	data := dataType{}
	result, err := ds.Get()
	c.Assert(err, check.IsNil)
	err = json.Unmarshal([]byte(result), &data)
	c.Assert(err, check.IsNil)
	c.Assert(data.Name, check.Equals, "Paul")
}

func (s *S) TestNew(c *check.C) {
	dsConfigTests := []struct {
		conf *DataSource
		err  error
	}{
		{&DataSource{URL: "http://tsuru.io", Method: "GET"}, nil},
		{&DataSource{URL: "http://tsuru.io"}, errors.New("datasource: method required")},
		{&DataSource{Method: ""}, errors.New("datasource: url required")},
	}
	for _, tt := range dsConfigTests {
		err := New(tt.conf)
		c.Check(err, check.DeepEquals, tt.err)
	}
}

func (s *S) TestGet(c *check.C) {
	ds := DataSource{
		Name:    "xpto",
		Headers: nil,
	}
	s.conn.DataSources().Insert(&ds)
	instance, err := Get(ds.Name)
	c.Assert(err, check.IsNil)
	c.Assert(instance, check.DeepEquals, &ds)

}
