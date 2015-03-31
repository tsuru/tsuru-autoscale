// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package datasource

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestRegister(c *check.C) {
	var ds dataSource
	dsFactory := func(conf map[string]interface{}) (dataSource, error) {
		return ds, nil
	}
	Register("graphite", dsFactory)
	d, err := New("graphite", nil)
	c.Assert(err, check.IsNil)
	c.Assert(d, check.DeepEquals, ds)
}

type testHandler struct{}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content := `{"name":"Paul"}`
	w.Write([]byte(content))
}

func (s *S) TestHttpDataSourceImplements(c *check.C) {
	ds := httpDataSource{}
	var expected dataSource
	c.Assert(&ds, check.Implements, &expected)
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

func (s *S) TestHttpDataSourceFactory(c *check.C) {
	dsConfigTests := []struct {
		conf map[string]interface{}
		err  error
	}{
		{nil, errors.New("datasource: url required")},
		{map[string]interface{}{"url": "", "method": "", "body": ""}, nil},
		{map[string]interface{}{"url": "", "body": ""}, errors.New("datasource: method required")},
		{map[string]interface{}{"url": "", "method": ""}, errors.New("datasource: body required")},
		{map[string]interface{}{"method": "", "body": ""}, errors.New("datasource: url required")},
	}
	for _, tt := range dsConfigTests {
		_, err := httpDataSourceFactory(tt.conf)
		c.Check(err, check.DeepEquals, tt.err)
	}
}

func (s *S) TestHttpDataSourceFactoryRegistered(c *check.C) {
	dsFactory, ok := dataSources["http"]
	c.Assert(ok, check.Equals, true)
	var expected dataSourceFactory
	c.Assert(dsFactory, check.FitsTypeOf, expected)
}

func (s *S) TestList(c *check.C) {
	var expected []string
	for name := range dataSources {
		expected = append(expected, name)
	}
	ds := List()
	c.Assert(ds, check.DeepEquals, expected)
}
