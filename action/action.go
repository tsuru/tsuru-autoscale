// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package action

import (
	"errors"
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
	return conn.DataSources().Insert(&a)
}

func (a *Action) Do() error {
	req, err := http.NewRequest(a.Method, a.URL, strings.NewReader(a.Body))
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
