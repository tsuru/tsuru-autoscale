// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tsuru/tsuru-autoscale/datasource"
	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru/db/dbtest"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	conn *db.Storage
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	var err error
	s.conn, err = db.Conn()
	c.Assert(err, check.IsNil)
}

func (s *S) TearDownTest(c *check.C) {
	dbtest.ClearAllCollections(s.conn.Actions().Database)
}

func (s *S) TestNewDataSource(c *check.C) {
	body := `{"name":"new","url":"http://tsuru.io","method":"GET"}`
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/datasource", strings.NewReader(body))
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestAllDataSources(c *check.C) {
	err := datasource.New(&datasource.DataSource{URL: "http://tsuru.io", Method: "GET"})
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/datasource", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var ds []datasource.DataSource
	err = json.Unmarshal(body, &ds)
	c.Assert(err, check.IsNil)
	c.Assert(ds, check.HasLen, 1)
}

func (s *S) TestRemoveDataSourceNotFound(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/datasource/notfound", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusNotFound)
}

func (s *S) TestRemoveDataSource(c *check.C) {
	ds := &datasource.DataSource{URL: "http://tsuru.io", Method: "GET", Name: "ds"}
	err := datasource.New(ds)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", fmt.Sprintf("/datasource/%s", ds.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	_, err = datasource.Get(ds.Name)
	c.Assert(err, check.NotNil)
}

func (s *S) TestGetDataSource(c *check.C) {
	ds := &datasource.DataSource{URL: "http://tsuru.io", Method: "GET", Name: "ds"}
	err := datasource.New(ds)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("/datasource/%s", ds.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var got datasource.DataSource
	err = json.Unmarshal(body, &got)
	c.Assert(err, check.IsNil)
	c.Assert(ds.Name, check.Equals, got.Name)
}
