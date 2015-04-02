// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	ID     bson.ObjectId `bson:"_id"`
	Name   string
	Params map[string]string
}

func serviceAdd(name string, params map[string]string) (*service, error) {
	srv := &service{
		ID:     bson.NewObjectId(),
		Name:   name,
		Params: params,
	}
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return srv, conn.Services().Insert(srv)
}
