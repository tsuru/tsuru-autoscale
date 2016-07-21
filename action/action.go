// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package action

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru-autoscale/log"
	"gopkg.in/mgo.v2/bson"
)

func logger() *log.Logger {
	return log.Log()
}

// Action represents an AutoScale action to increase or decrease the
// number of the units.
type Action struct {
	Name    string
	URL     string
	Method  string
	Body    string
	Headers map[string]string
}

// New creates a new action.
func New(a *Action) error {
	if a.URL == "" {
		return errors.New("action: url required")
	}
	if a.Method == "" {
		return errors.New("action: method required")
	}
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return err
	}
	defer conn.Close()
	return conn.Actions().Insert(&a)
}

// FindByName finds action by name.
func FindByName(name string) (*Action, error) {
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer conn.Close()
	var action Action
	err = conn.Actions().Find(bson.M{"name": name}).One(&action)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	return &action, nil
}

// Remove removes an action.
func Remove(a *Action) error {
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return err
	}
	defer conn.Close()
	return conn.Actions().Remove(bson.M{"name": a.Name})
}

// All return a list of all actions
func All() ([]Action, error) {
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer conn.Close()
	var actions []Action
	err = conn.Actions().Find(nil).All(&actions)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	return actions, nil
}

// Do executes the action
func (a *Action) Do(appName string, envs map[string]string) error {
	body := a.Body
	url := strings.Replace(a.URL, "{app}", appName, -1)
	for key, value := range envs {
		body = strings.Replace(body, fmt.Sprintf("{%s}", key), value, -1)
		url = strings.Replace(url, fmt.Sprintf("{%s}", key), value, -1)
	}
	logger().Printf("action %s - url: %s - body: %s - method: %s", a.Name, url, body, a.Method)
	req, err := http.NewRequest(a.Method, url, strings.NewReader(body))
	if err != nil {
		logger().Error(err)
		return err
	}
	for key, value := range a.Headers {
		req.Header.Add(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger().Error(err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
