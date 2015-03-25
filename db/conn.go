// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package db encapsulates tsuru-autoscale connection with MongoDB.
//
// The function Conn dials to MongoDB using data from the configuration file
// and returns a connection (represented by the storage.Storage type). It
// manages an internal pool of connections, and reconnects in case of failures.
// That means that you should not store references to the connection, but
// always call Open.
package db

import (
	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/db/storage"
)

const (
	DefaultDatabaseURL  = "127.0.0.1:27017"
	DefaultDatabaseName = "tsuru_autoscale"
)

type Storage struct {
	*storage.Storage
}

// conn reads the tsuru-autoscale config and calls storage.Open to get a database connection.
func conn() (*storage.Storage, error) {
	url, _ := config.GetString("database:url")
	if url == "" {
		url = DefaultDatabaseURL
	}
	dbname, _ := config.GetString("database:name")
	if dbname == "" {
		dbname = DefaultDatabaseName
	}
	return storage.Open(url, dbname)
}

func Conn() (*Storage, error) {
	var (
		strg Storage
		err  error
	)
	strg.Storage, err = conn()
	return &strg, err
}

// Event returns the event collection from MongoDB.
func (s *Storage) Event() *storage.Collection {
	c := s.Collection("event")
	return c
}

// Config returns the config collection from MongoDB.
func (s *Storage) Config() *storage.Collection {
	c := s.Collection("config")
	return c
}
