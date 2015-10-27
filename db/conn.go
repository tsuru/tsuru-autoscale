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
	"os"

	"github.com/tsuru/tsuru/db/storage"
	"gopkg.in/mgo.v2"
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
	url := os.Getenv("MONGODB_URL")
	if url == "" {
		url = DefaultDatabaseURL
	}
	dbname := os.Getenv("MONGODB_DATABASE_NAME")
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

// Events returns the events collection from MongoDB.
func (s *Storage) Events() *storage.Collection {
	c := s.Collection("events")
	return c
}

// Configs returns the configs collection from MongoDB.
func (s *Storage) Configs() *storage.Collection {
	c := s.Collection("configs")
	return c
}

// Instances returns the instances collection from MongoDB.
func (s *Storage) Instances() *storage.Collection {
	nameIndex := mgo.Index{Key: []string{"name"}, Unique: true}
	c := s.Collection("instances")
	c.EnsureIndex(nameIndex)
	return c
}

// DataSources returns the datasources collection from MongoDB.
func (s *Storage) DataSources() *storage.Collection {
	nameIndex := mgo.Index{Key: []string{"name"}, Unique: true}
	c := s.Collection("datasources")
	c.EnsureIndex(nameIndex)
	return c
}

// Alarms returns the alarms collection from MongoDB.
func (s *Storage) Alarms() *storage.Collection {
	nameIndex := mgo.Index{Key: []string{"name"}, Unique: true}
	c := s.Collection("alarms")
	c.EnsureIndex(nameIndex)
	return c
}

// Actions returns the actions collection from MongoDB.
func (s *Storage) Actions() *storage.Collection {
	nameIndex := mgo.Index{Key: []string{"name"}, Unique: true}
	c := s.Collection("actions")
	c.EnsureIndex(nameIndex)
	return c
}

// Wizard returns the wizard collection from MongoDB.
func (s *Storage) Wizard() *storage.Collection {
	nameIndex := mgo.Index{Key: []string{"name"}, Unique: true}
	c := s.Collection("wizard")
	c.EnsureIndex(nameIndex)
	return c
}
