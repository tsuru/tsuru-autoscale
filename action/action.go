// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package action

import "net/url"

// Action represents an AutoScale action to increase or decrease the
// number of the units.
type Action struct {
	Name string
	URL  *url.URL
}

func New(name string, url *url.URL) (*Action, error) {
	return &Action{Name: name, URL: url}, nil
}
