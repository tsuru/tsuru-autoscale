// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func init() {
	Register("http", httpDataSourceFactory)
}

// dataSource represents a data source.
type dataSource interface {
	// Get gets data from data source and
	// parses the JSON-encoded data and stores the result
	// in the value pointed to by v.
	Get(v interface{}) error
}

type dataSourceFactory func(conf map[string]interface{}) (dataSource, error)

var dataSources = make(map[string]dataSourceFactory)

// Register registers a new dataSource.
func Register(name string, ds dataSourceFactory) {
	dataSources[name] = ds
}

// NewDataSource creates a new data source instance.
func NewDataSource(name string, conf map[string]interface{}) (dataSource, error) {
	return dataSources[name](conf)
}

type httpDataSource struct {
	url    string
	method string
	body   string
}

func httpDataSourceFactory(conf map[string]interface{}) (dataSource, error) {
	url, ok := conf["url"]
	if !ok {
		return nil, errors.New("url required")
	}
	method, ok := conf["method"]
	if !ok {
		return nil, errors.New("method required")
	}
	body, ok := conf["body"]
	if !ok {
		return nil, errors.New("body required")
	}
	ds := httpDataSource{
		url:    url.(string),
		method: method.(string),
		body:   body.(string),
	}
	return &ds, nil
}

func (ds *httpDataSource) Get(v interface{}) error {
	req, err := http.NewRequest(ds.method, ds.url, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
