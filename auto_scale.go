// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func StartAutoScale() {
	go runAutoScale()
}

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
	return &evt, conn.AutoScale().Insert(evt)
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
	return conn.AutoScale().UpdateId(evt.ID, evt)
}

// Action represents an AutoScale action to increase or decrease the
// number of the units.
type Action struct {
	Wait       time.Duration `json:"wait"`
	Expression string        `json:"expression"`
	Units      uint          `json:"units"`
}

func NewAction(expression string, units uint, wait time.Duration) (*Action, error) {
	if expressionIsValid(expression) {
		return &Action{Wait: wait, Expression: expression, Units: units}, nil
	}
	return nil, errors.New("Expression is not valid.")
}

var expressionRegex = regexp.MustCompile("{(.*)} ([><=]) ([0-9]+)")

func expressionIsValid(expression string) bool {
	return expressionRegex.MatchString(expression)
}

func (action *Action) metric() string {
	return expressionRegex.FindStringSubmatch(action.Expression)[1]
}

func (action *Action) operator() string {
	return expressionRegex.FindStringSubmatch(action.Expression)[2]
}

func (action *Action) value() (float64, error) {
	return strconv.ParseFloat(expressionRegex.FindStringSubmatch(action.Expression)[3], 64)
}

// Config represents the configuration for the auto scale.
type Config struct {
	Name     string `json:"increase"`
	Increase Action `json:"increase"`
	Decrease Action `json:"decrease"`
	MinUnits uint   `json:"minUnits"`
	MaxUnits uint   `json:"maxUnits"`
	Enabled  bool   `json:"enabled"`
}

type App struct {
	Name string
}

func (a *App) Units() []string {
	return nil
}

func (a *App) Metric(kind string) (float64, error) {
	return float64(0), nil
}

func (a *App) AddUnits(n uint, writer io.Writer) error {
	return nil
}

func (a *App) RemoveUnits(n uint) error {
	return nil
}

func runAutoScaleOnce() {
	configs := []Config{}
	for _, config := range configs {
		err := scaleIfNeeded(&config)
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func runAutoScale() {
	for {
		runAutoScaleOnce()
		time.Sleep(30 * time.Second)
	}
}

func scaleIfNeeded(config *Config) error {
	if config == nil {
		return errors.New("AutoScale is not configured.")
	}
	/*
		increaseMetric, _ := app.Metric(config.Increase.metric())
		value, _ := config.Increase.value()
		if increaseMetric > value {
			currentUnits := uint(len(app.Units()))
			maxUnits := config.MaxUnits
			if maxUnits == 0 {
				maxUnits = 1
			}
			if currentUnits >= maxUnits {
				return nil
			}
			if wait, err := shouldWait(app, config.Increase.Wait); err != nil {
				return err
			} else if wait {
				return nil
			}
			evt, err := NewEvent(app, "increase")
			if err != nil {
				return fmt.Errorf("Error trying to insert auto scale event, auto scale aborted: %s", err.Error())
		 	}
			inc := config.Increase.Units
			if currentUnits+inc > config.MaxUnits {
				inc = config.MaxUnits - currentUnits
			}
			addUnitsErr := app.AddUnits(inc, nil)
			err = evt.update(addUnitsErr)
			if err != nil {
				log.Errorf("Error trying to update auto scale event: %s", err.Error())
			}
			return addUnitsErr
		}
		decreaseMetric, _ := app.Metric(config.Decrease.metric())
		value, _ = config.Decrease.value()
		if decreaseMetric < value {
			currentUnits := uint(len(app.Units()))
			minUnits := config.MinUnits
			if minUnits == 0 {
				minUnits = 1
			}
			if currentUnits <= minUnits {
				return nil
			}
			if wait, err := shouldWait(app, config.Decrease.Wait); err != nil {
				return err
			} else if wait {
				return nil
			}
			evt, err := NewEvent(app, "decrease")
			if err != nil {
				return fmt.Errorf("Error trying to insert auto scale event, auto scale aborted: %s", err.Error())
			}
			dec := config.Decrease.Units
			if currentUnits-dec < config.MinUnits {
				dec = currentUnits - config.MinUnits
			}
			removeUnitsErr := app.RemoveUnits(dec)
			err = evt.update(removeUnitsErr)
			if err != nil {
				log.Errorf("Error trying to update auto scale event: %s", err.Error())
			}
			return removeUnitsErr
		}
	*/
	return nil
}

func shouldWait(app *App, waitPeriod time.Duration) (bool, error) {
	now := time.Now().UTC()
	lastEvent, err := lastScaleEvent(app.Name)
	if err != nil && err != mgo.ErrNotFound {
		return false, err
	}
	if err != mgo.ErrNotFound && lastEvent.EndTime.IsZero() {
		return true, nil
	}
	diff := now.Sub(lastEvent.EndTime)
	if diff > waitPeriod {
		return false, nil
	}
	return true, nil
}

func lastScaleEvent(appName string) (Event, error) {
	var event Event
	conn, err := db.Conn()
	if err != nil {
		return event, err
	}
	defer conn.Close()
	err = conn.AutoScale().Find(bson.M{"appname": appName}).Sort("-starttime").One(&event)
	return event, err
}

func ListAutoScaleHistory(appName string) ([]Event, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var history []Event
	q := bson.M{}
	if appName != "" {
		q["appname"] = appName
	}
	err = conn.AutoScale().Find(q).Sort("-_id").Limit(200).All(&history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func AutoScaleEnable(config *Config) error {
	config.Enabled = true
	return nil
}

func AutoScaleDisable(config *Config) error {
	config.Enabled = false
	return nil
}
