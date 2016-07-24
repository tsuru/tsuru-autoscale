// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alarm

import (
	"time"

	"github.com/tsuru/tsuru-autoscale/action"
	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/mgo.v2/bson"
)

// Event represents an auto scale event with
// the scale metadata.
type Event struct {
	ID         bson.ObjectId `bson:"_id"`
	StartTime  time.Time
	EndTime    time.Time `bson:",omitempty"`
	Alarm      *Alarm
	Successful bool
	Error      string `bson:",omitempty"`
	Action     *action.Action
}

// NewEvent creates a new alarm event
func NewEvent(alarm *Alarm, action *action.Action) (*Event, error) {
	evt := Event{
		ID:        bson.NewObjectId(),
		StartTime: time.Now().UTC(),
		Alarm:     alarm,
		Action:    action,
	}
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
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
		logger().Error(err)
		return err
	}
	defer conn.Close()
	return conn.Events().UpdateId(evt.ID, evt)
}

func lastScaleEvent(alarm *Alarm) (Event, error) {
	var event Event
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return event, err
	}
	defer conn.Close()
	err = conn.Events().Find(bson.M{"alarm.name": alarm.Name}).Sort("-starttime").One(&event)
	return event, err
}

// FindEventsBy is an extensible way to find events by query
func FindEventsBy(q bson.M, limit int) ([]Event, error) {
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer conn.Close()
	var events []Event
	err = conn.Events().Find(q).Sort("-starttime").Limit(limit).All(&events)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	return events, nil
}

// EventsByAlarmName returns a list of events by alarm name
func EventsByAlarmName(alarm string) ([]Event, error) {
	q := bson.M{}
	if alarm != "" {
		q["alarm.name"] = alarm
	}
	return FindEventsBy(q, 200)
}
