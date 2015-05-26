// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package action

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru/db/dbtest"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	conn *db.Storage
}

func (s *S) SetUpSuite(c *check.C) {
	err := os.Setenv("MONGODB_DATABASE_NAME", "tsuru_autoscale_action")
	c.Assert(err, check.IsNil)
	s.conn, err = db.Conn()
	c.Assert(err, check.IsNil)
}

func (s *S) TearDownTest(c *check.C) {
	dbtest.ClearAllCollections(s.conn.Actions().Database)
}

func (s *S) TearDownSuite(c *check.C) {
	err := os.Unsetenv("MONGODB_DATABASE_NAME")
	c.Assert(err, check.IsNil)
}

var _ = check.Suite(&S{})

func (s *S) TestNew(c *check.C) {
	actionTests := []struct {
		a   *Action
		err error
	}{
		{&Action{URL: "http://tsuru.io", Method: "GET"}, nil},
		{&Action{URL: "http://tsuru.io"}, errors.New("action: method required")},
		{&Action{Method: ""}, errors.New("action: url required")},
	}
	for _, tt := range actionTests {
		err := New(tt.a)
		c.Check(err, check.DeepEquals, tt.err)
	}
}

func (s *S) TestDo(c *check.C) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		c.Assert(string(body), check.Equals, `{"units": "1"}`)
	}))
	defer ts.Close()
	a := Action{URL: ts.URL, Method: "GET", Body: `{"units": "{step}"}`}
	err := New(&a)
	c.Assert(err, check.IsNil)
	envs := map[string]string{"step": "1"}
	err = a.Do("app", envs)
	c.Assert(err, check.IsNil)
	c.Assert(called, check.Equals, true)
}

func (s *S) TestAll(c *check.C) {
	a := Action{
		Name:    "xpto",
		Headers: nil,
	}
	s.conn.Actions().Insert(&a)
	a = Action{
		Name:    "xpto2",
		Headers: nil,
	}
	s.conn.Actions().Insert(&a)
	all, err := All()
	c.Assert(err, check.IsNil)
	c.Assert(all, check.HasLen, 2)
}

func (s *S) TestFindByName(c *check.C) {
	a := Action{
		Name:    "xpto123",
		Headers: map[string]string{},
	}
	s.conn.Actions().Insert(&a)
	a = Action{
		Name:    "xpto1234",
		Headers: map[string]string{},
	}
	s.conn.Actions().Insert(&a)
	na, err := FindByName(a.Name)
	c.Assert(err, check.IsNil)
	c.Assert(na, check.DeepEquals, &a)
}

func (s *S) TestRemove(c *check.C) {
	a := Action{
		Name:    "xpto123",
		Headers: map[string]string{},
	}
	s.conn.Actions().Insert(&a)
	err := Remove(&a)
	c.Assert(err, check.IsNil)
	_, err = FindByName(a.Name)
	c.Assert(err, check.NotNil)
}
