// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestServiceAdd(c *check.C) {
	name := "name"
	params := map[string]string{
		"key": "value",
	}
	srv, err := serviceAdd(name, params)
	c.Assert(err, check.IsNil)
	c.Assert(srv.Name, check.Equals, name)
	c.Assert(srv.Params, check.DeepEquals, params)
}
