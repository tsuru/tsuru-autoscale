// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestAutoScale(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	autoscale := strg.AutoScale()
	autoscalec := strg.Collection("autoscale")
	c.Assert(autoscale, check.DeepEquals, autoscalec)
}
