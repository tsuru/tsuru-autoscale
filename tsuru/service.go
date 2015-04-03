// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/mgo.v2/bson"
)

type instance struct {
	ID       bson.ObjectId `bson:"_id"`
	Name     string
	Metadata map[string]string
	Apps     []string
}

func NewInstance(name string, metadata map[string]string) (*instance, error) {
	i := &instance{
		ID:       bson.NewObjectId(),
		Name:     name,
		Metadata: metadata,
	}
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return i, conn.Instances().Insert(i)
}
