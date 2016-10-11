// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	err := os.Setenv("MONGODB_DATABASE_NAME", "tsuru_autoscale_api")
	c.Assert(err, check.IsNil)
	s.conn, err = db.Conn()
	c.Assert(err, check.IsNil)
}

func (s *S) TearDownTest(c *check.C) {
	dbtest.ClearAllCollections(s.conn.Actions().Database)
	dbtest.ClearAllCollections(s.conn.Alarms().Database)
	dbtest.ClearAllCollections(s.conn.DataSources().Database)
	dbtest.ClearAllCollections(s.conn.Events().Database)
	dbtest.ClearAllCollections(s.conn.Instances().Database)
	dbtest.ClearAllCollections(s.conn.Wizard().Database)
}

func (s *S) TestNewDataSource(c *check.C) {
	body := `{"name":"new","url":"http://tsuru.io","method":"GET"}`
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/datasource", strings.NewReader(body))
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestFindByDataSources(c *check.C) {
	err := datasource.New(&datasource.DataSource{
		URL:    "http://tsuru.io",
		Method: "GET",
		Public: false,
	})
	c.Assert(err, check.IsNil)
	var tests = []struct {
		url    string
		length int
	}{
		{"/datasource", 1},
		{"/datasource?public=true", 0},
	}
	for _, t := range tests {
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", t.url, nil)
		c.Check(err, check.IsNil)
		request.Header.Add("Authorization", "token")
		server(recorder, request)
		c.Check(recorder.Code, check.Equals, http.StatusOK)
		c.Check(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
		body := recorder.Body.Bytes()
		var ds []datasource.DataSource
		err = json.Unmarshal(body, &ds)
		c.Check(err, check.IsNil)
		c.Check(ds, check.HasLen, t.length)
	}
}

func (s *S) TestRemoveDataSourceNotFound(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/datasource/notfound", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
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
	server(recorder, request)
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
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var got datasource.DataSource
	err = json.Unmarshal(body, &got)
	c.Assert(err, check.IsNil)
	c.Assert(ds.Name, check.Equals, got.Name)
}
