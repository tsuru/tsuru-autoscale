// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"gopkg.in/check.v1"
)

func (s *S) TestActionMetric(c *check.C) {
	a := &Action{Expression: "{cpu} > 80"}
	c.Assert(a.metric(), check.Equals, "cpu")
}

func (s *S) TestActionOperator(c *check.C) {
	a := &Action{Expression: "{cpu} > 80"}
	c.Assert(a.operator(), check.Equals, ">")
}

func (s *S) TestActionValue(c *check.C) {
	a := &Action{Expression: "{cpu} > 80"}
	value, err := a.value()
	c.Assert(err, check.IsNil)
	c.Assert(value, check.Equals, float64(80))
}

func (s *S) TestValidateExpression(c *check.C) {
	cases := map[string]bool{
		"{cpu} > 10": true,
		"{cpu} = 10": true,
		"{cpu} < 10": true,
		"cpu < 10":   false,
		"{cpu} 10":   false,
		"{cpu} <":    false,
		"{cpu}":      false,
		"<":          false,
		"100":        false,
	}
	for expression, expected := range cases {
		c.Assert(expressionIsValid(expression), check.Equals, expected)
	}
}

func (s *S) TestNewAction(c *check.C) {
	expression := "{cpu} > 10"
	units := uint(2)
	wait := time.Second
	a, err := NewAction(expression, units, wait)
	c.Assert(err, check.IsNil)
	c.Assert(a.Expression, check.Equals, expression)
	c.Assert(a.Units, check.Equals, units)
	c.Assert(a.Wait, check.Equals, wait)
	expression = "{cpu} >"
	units = uint(2)
	wait = time.Second
	a, err = NewAction(expression, units, wait)
	c.Assert(err, check.NotNil)
	c.Assert(a, check.IsNil)
}
