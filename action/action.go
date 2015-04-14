// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package action

import (
	"net/http"
	"net/url"
	"strings"
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

func New(name string, url *url.URL) (*Action, error) {
	return &Action{Name: name, URL: url.String()}, nil
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
