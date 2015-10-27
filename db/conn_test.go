// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"reflect"
	"testing"

	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

type hasUniqueIndexChecker struct{}

func (c *hasUniqueIndexChecker) Info() *check.CheckerInfo {
	return &check.CheckerInfo{Name: "HasUniqueField", Params: []string{"collection", "key"}}
}

func (c *hasUniqueIndexChecker) Check(params []interface{}, names []string) (bool, string) {
	collection, ok := params[0].(*storage.Collection)
	if !ok {
		return false, "first parameter should be a Collection"
	}
	key, ok := params[1].([]string)
	if !ok {
		return false, "second parameter should be the key, as used for mgo index declaration (slice of strings)"
	}
	indexes, err := collection.Indexes()
	if err != nil {
		return false, "failed to get collection indexes: " + err.Error()
	}
	for _, index := range indexes {
		if reflect.DeepEqual(index.Key, key) {
			return index.Unique, ""
		}
	}
	return false, ""
}

var HasUniqueIndex check.Checker = &hasUniqueIndexChecker{}

func (s *S) TestEvents(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	event := strg.Events()
	eventc := strg.Collection("events")
	c.Assert(event, check.DeepEquals, eventc)
}

func (s *S) TestConfigs(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	config := strg.Configs()
	configc := strg.Collection("configs")
	c.Assert(config, check.DeepEquals, configc)
}

func (s *S) TestInstances(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	instance := strg.Instances()
	instancec := strg.Collection("instances")
	c.Assert(instance, check.DeepEquals, instancec)
	c.Assert(instance, HasUniqueIndex, []string{"name"})
}

func (s *S) TestDataSources(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	datasource := strg.DataSources()
	datasourcec := strg.Collection("datasources")
	c.Assert(datasource, check.DeepEquals, datasourcec)
	c.Assert(datasource, HasUniqueIndex, []string{"name"})
}

func (s *S) TestAlarms(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	alarm := strg.Alarms()
	alarmc := strg.Collection("alarms")
	c.Assert(alarm, check.DeepEquals, alarmc)
	c.Assert(alarm, HasUniqueIndex, []string{"name"})
}

func (s *S) TestActions(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	action := strg.Actions()
	actionc := strg.Collection("actions")
	c.Assert(action, check.DeepEquals, actionc)
	c.Assert(action, HasUniqueIndex, []string{"name"})
}

func (s *S) TestWizard(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	wizard := strg.Wizard()
	wizardc := strg.Collection("wizard")
	c.Assert(wizard, check.DeepEquals, wizardc)
	c.Assert(wizard, HasUniqueIndex, []string{"name"})
}
