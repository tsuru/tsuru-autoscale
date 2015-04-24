// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import "github.com/tsuru/tsuru-autoscale/db"

// Instance represents a tsuru service instance.
type Instance struct {
	Name string
	User string
	Team string
	Apps []string
}

func NewInstance(i *Instance) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Instances().Insert(i)
}
