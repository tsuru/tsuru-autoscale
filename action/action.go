// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package action

import (
	"io"
	"net/http"
	"net/url"
)

// Action represents an AutoScale action to increase or decrease the
// number of the units.
type Action struct {
	Name string
	URL  *url.URL
	Method string
	Body io.Reader
}

func New(name string, url *url.URL) (*Action, error) {
	return &Action{Name: name, URL: url}, nil
}

func (a *Action) Do() error {
	req, err := http.NewRequest(a.Method, a.URL.String(), a.Body)
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	return err
}
