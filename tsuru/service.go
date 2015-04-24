// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/mgo.v2/bson"
)

// Instance represents a tsuru service instance.
type Instance struct {
	Name string
	User string
	Team string
	Apps []string
}

// NewInstance creates a new service instance.
func NewInstance(i *Instance) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Instances().Insert(i)
}

// GetInstanceByName finds a service instance by name.
func GetInstanceByName(name string) (*Instance, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var i Instance
	err = conn.Instances().Find(bson.M{"name": name}).One(&i)
	if err != nil {
		return nil, err
	}
	return &i, nil
}
