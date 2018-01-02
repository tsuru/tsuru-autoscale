// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru-autoscale/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func logger() *log.Logger {
	return log.Log()
}

// Instance represents a tsuru service instance.
type Instance struct {
	ID   bson.ObjectId `bson:"_id" json:"-"`
	Name string
	User string
	Team string
	Apps []string `json:",omitempty"`
}

func (i *Instance) update() error {
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
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
func (i *Instance) AddApp(app, host string) error {
	if app == "" {
		app = appFromHost(host)
	}
	if contains(i.Apps, app) {
		return errors.New("")
	}
	i.Apps = append(i.Apps, app)
	return i.update()
}

// RemoveApp removes app from an instance.
func (i *Instance) RemoveApp(app, host string) error {
	if app == "" {
		app = appFromHost(host)
	}
	if !contains(i.Apps, app) {
		return errors.New("")
	}
	var apps []string
	for _, a := range i.Apps {
		if a != app {
			apps = append(apps, a)
		}
	}
	i.Apps = apps
	return i.update()
}

// NewInstance creates a new service instance.
func NewInstance(i *Instance) error {
	if i.ID.Hex() == "" {
		i.ID = bson.NewObjectId()
	}
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return err
	}
	defer conn.Close()
	return conn.Instances().Insert(i)
}

// RemoveInstance removes an auto scale instance
func RemoveInstance(i *Instance) error {
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return err
	}
	defer conn.Close()
	return conn.Instances().Remove(bson.M{"name": i.Name})
}

// GetInstanceByName finds a service instance by name.
func GetInstanceByName(name string) (*Instance, error) {
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer conn.Close()
	var i Instance
	err = conn.Instances().Find(bson.M{"name": name}).One(&i)
	if err != nil {
		if err == mgo.ErrNotFound {
			err = fmt.Errorf("instance %q not found", name)
		}
		logger().Error(err)
		return nil, err
	}
	return &i, nil
}
