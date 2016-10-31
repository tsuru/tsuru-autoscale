// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/ajg/form"
	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/datasource"
	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru/db/dbtest"
	"gopkg.in/check.v1"
)

func server(w http.ResponseWriter, r *http.Request) {
	m := mux.NewRouter()
	Router(m)
	m.ServeHTTP(w, r)
}

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	conn *db.Storage
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	err := os.Setenv("MONGODB_DATABASE_NAME", "tsuru_autoscale_web")
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

func (s *S) TestDataSourceAdd(c *check.C) {
	recorder := httptest.NewRecorder()
	ds := datasource.DataSource{
		Name:   "new",
		URL:    "http://tsuru.io",
		Method: "GET",
	}
	v, err := form.EncodeToValues(&ds)
	c.Assert(err, check.IsNil)
	body := strings.NewReader(v.Encode())
	request, err := http.NewRequest("POST", "/datasource/add", body)
	c.Assert(err, check.IsNil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
}

func (s *S) TestDataSourceAddEmptyHeader(c *check.C) {
	recorder := httptest.NewRecorder()
	ds := datasource.DataSource{
		Name:    "new",
		URL:     "http://tsuru.io",
		Method:  "GET",
		Headers: map[string]string{" ": " "},
	}
	v := url.Values{
		"key":    []string{"", "f", ""},
		"value":  []string{"", "f", ""},
		"name":   []string{"new"},
		"url":    []string{"sdfasd"},
		"method": []string{"GET"},
	}
	body := strings.NewReader(v.Encode())
	request, err := http.NewRequest("POST", "/datasource/add", body)
	c.Assert(err, check.IsNil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
	r, err := datasource.Get(ds.Name)
	c.Assert(err, check.IsNil)
	c.Assert(len(r.Headers), check.Equals, 1)
}

func (s *S) TestDataSourceRemove(c *check.C) {
	recorder := httptest.NewRecorder()
	ds := datasource.DataSource{
		Name:   "new",
		URL:    "http://tsuru.io",
		Method: "GET",
	}
	err := datasource.New(&ds)
	c.Assert(err, check.IsNil)
	request, err := http.NewRequest("GET", "/datasource/new/delete", nil)
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
	_, err = datasource.Get(ds.Name)
	c.Assert(err, check.NotNil)
}
