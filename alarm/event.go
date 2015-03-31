// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alarm

import (
	"time"

	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/mgo.v2/bson"
)

// Event represents an auto scale event with
// the scale metadata.
type Event struct {
	ID         bson.ObjectId `bson:"_id"`
	StartTime  time.Time
	EndTime    time.Time `bson:",omitempty"`
	Config     *Config
	Type       string
	Successful bool
	Error      string `bson:",omitempty"`
}

func NewEvent(config *Config, scaleType string) (*Event, error) {
	evt := Event{
		ID:        bson.NewObjectId(),
		StartTime: time.Now().UTC(),
		Config:    config,
		Type:      scaleType,
	}
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return &evt, conn.Events().Insert(evt)
}

func (evt *Event) update(err error) error {
	if err != nil {
		evt.Error = err.Error()
	}
	evt.Successful = err == nil
	evt.EndTime = time.Now().UTC()
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Events().UpdateId(evt.ID, evt)
}

func lastScaleEvent(config *Config) (Event, error) {
	var event Event
	conn, err := db.Conn()
	if err != nil {
		return event, err
	}
	defer conn.Close()
	err = conn.Events().Find(bson.M{"config.name": config.Name}).Sort("-starttime").One(&event)
	return event, err
}

func eventsByConfigName(config *Config) ([]Event, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var events []Event
	q := bson.M{}
	if config != nil {
		q["config.name"] = config.Name
	}
	err = conn.Events().Find(q).Sort("-_id").Limit(200).All(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
