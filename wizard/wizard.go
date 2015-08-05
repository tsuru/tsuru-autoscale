// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wizard

type autoscale struct {
	name      string
	scaleUp   scaleAction
	scaleDown scaleAction
	minUnits  int
}

type scaleAction struct {
	metric   string
	operator string
	value    string
	step     int
	waitTime int
}
