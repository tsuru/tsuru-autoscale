// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"errors"
	"strings"

	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/mgo.v2/bson"
)

// Instance represents a tsuru service instance.
type Instance struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string
	User string
	Team string
	Apps []string
}

func (i *Instance) update() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Instances().UpdateId(i.ID, i)
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func appFromHost(host string) string {
	return strings.Split(host, ".")[0]
}

// AddApp add new app to an instance.
func (i *Instance) AddApp(host string) error {
	app := appFromHost(host)
	if contains(i.Apps, app) {
		return errors.New("")
	}
	i.Apps = append(i.Apps, app)
	return i.update()
}

// NewInstance creates a new service instance.
func NewInstance(i *Instance) error {
	if i.ID.Hex() == "" {
		i.ID = bson.NewObjectId()
	}
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
