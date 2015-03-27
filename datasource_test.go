// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"gopkg.in/check.v1"
)

func (s *S) TestRegister(c *check.C) {
	var ds dataSource
	dsFactory := func(conf map[string]interface{}) (dataSource, error) {
		return ds, nil
	}
	Register("graphite", dsFactory)
	d, err := dataSources["graphite"](nil)
	c.Assert(err, check.IsNil)
	c.Assert(d, check.DeepEquals, ds)
}
