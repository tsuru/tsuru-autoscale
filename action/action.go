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
)

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
		return err
	}
	defer conn.Close()
	return conn.Actions().Insert(&a)
}

// FindByName finds action by name.
func FindByName(name string) (*Action, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var action Action
	err = conn.Actions().Find(nil).One(&action)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// Remove removes an action.
func Remove(a *Action) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Actions().Remove(a)
}

func All() ([]Action, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var actions []Action
	err = conn.Actions().Find(nil).All(&actions)
	if err != nil {
		return nil, err
	}
	return actions, nil
}

func (a *Action) Do(envs map[string]string) error {
	body := a.Body
	for key, value := range envs {
		body = strings.Replace(body, fmt.Sprintf("{%s}", key), value, -1)
	}
	req, err := http.NewRequest(a.Method, a.URL, strings.NewReader(body))
	if err != nil {
		return err
	}
	for key, value := range a.Headers {
		req.Header.Add(key, value)
	}
	client := &http.Client{}
	_, err = client.Do(req)
	return err
}
