// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
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
	ID              bson.ObjectId `bson:"_id"`
	AppName         string
	StartTime       time.Time
	EndTime         time.Time `bson:",omitempty"`
	Config		*Config
	Type            string
	Successful      bool
	Error           string `bson:",omitempty"`
}

func NewEvent(a *App, scaleType string) (*Event, error) {
	evt := Event{
		ID:              bson.NewObjectId(),
		StartTime:       time.Now().UTC(),
		Config: a.Config,
		AppName:         a.Name,
		Type:            scaleType,
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
	Name	 string `json:"increase"`
	Increase Action `json:"increase"`
	Decrease Action `json:"decrease"`
	MinUnits uint   `json:"minUnits"`
	MaxUnits uint   `json:"maxUnits"`
	Enabled  bool   `json:"enabled"`
}

type App struct {
	Config *Config
	Name            string
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

func autoScalableApps() ([]App, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var apps []App
	return apps, nil
}

func runAutoScaleOnce() {
	apps, err := autoScalableApps()
	if err != nil {
		log.Error(err.Error())
	}
	for _, app := range apps {
		err := scaleApplicationIfNeeded(&app)
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

func scaleApplicationIfNeeded(app *App) error {
	if app.Config == nil {
		return errors.New("AutoScale is not configured.")
	}
	increaseMetric, _ := app.Metric(app.Config.Increase.metric())
	value, _ := app.Config.Increase.value()
	if increaseMetric > value {
		currentUnits := uint(len(app.Units()))
		maxUnits := app.Config.MaxUnits
		if maxUnits == 0 {
			maxUnits = 1
		}
		if currentUnits >= maxUnits {
			return nil
		}
		if wait, err := shouldWait(app, app.Config.Increase.Wait); err != nil {
			return err
		} else if wait {
			return nil
		}
		evt, err := NewEvent(app, "increase")
		if err != nil {
			return fmt.Errorf("Error trying to insert auto scale event, auto scale aborted: %s", err.Error())
		}
		inc := app.Config.Increase.Units
		if currentUnits+inc > app.Config.MaxUnits {
			inc = app.Config.MaxUnits - currentUnits
		}
		addUnitsErr := app.AddUnits(inc, nil)
		err = evt.update(addUnitsErr)
		if err != nil {
			log.Errorf("Error trying to update auto scale event: %s", err.Error())
		}
		return addUnitsErr
	}
	decreaseMetric, _ := app.Metric(app.Config.Decrease.metric())
	value, _ = app.Config.Decrease.value()
	if decreaseMetric < value {
		currentUnits := uint(len(app.Units()))
		minUnits := app.Config.MinUnits
		if minUnits == 0 {
			minUnits = 1
		}
		if currentUnits <= minUnits {
			return nil
		}
		if wait, err := shouldWait(app, app.Config.Decrease.Wait); err != nil {
			return err
		} else if wait {
			return nil
		}
		evt, err := NewEvent(app, "decrease")
		if err != nil {
			return fmt.Errorf("Error trying to insert auto scale event, auto scale aborted: %s", err.Error())
		}
		dec := app.Config.Decrease.Units
		if currentUnits-dec < app.Config.MinUnits {
			dec = currentUnits - app.Config.MinUnits
		}
		removeUnitsErr := app.RemoveUnits(dec)
		err = evt.update(removeUnitsErr)
		if err != nil {
			log.Errorf("Error trying to update auto scale event: %s", err.Error())
		}
		return removeUnitsErr
	}
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

func AutoScaleEnable(app *App) error {
	if app.Config == nil {
		app.Config = &Config{}
	}
	app.Config.Enabled = true
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	return nil
}

func AutoScaleDisable(app *App) error {
	if app.Config == nil {
		app.Config = &Config{}
	}
	app.Config.Enabled = false
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	return nil
}

func SetConfig(app *App, config *Config) error {
	app.Config = config
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	return nil
}
