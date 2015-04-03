// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package datasource

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/mgo.v2/bson"
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

type dataSourceFactory func(metadata map[string]string) (dataSource, error)

var dataSources = make(map[string]dataSourceFactory)

// Register registers a new dataSource.
func Register(name string, ds dataSourceFactory) {
	dataSources[name] = ds
}

type Instance struct {
	Name     string
	Metadata map[string]string
}

// New creates a new data source instance.
func New(name string, metadata map[string]string) (dataSource, error) {
	instance := Instance{
		Name:     name,
		Metadata: metadata,
	}
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	err = conn.DataSources().Insert(&instance)
	if err != nil {
		return nil, err
	}
	return dataSources[name](metadata)
}

func Get(name string) (*Instance, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var i Instance
	err = conn.Instances().Find(bson.M{"name": name}).One(&i)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

type httpDataSource struct {
	url    string
	method string
	body   string
}

func httpDataSourceFactory(metadata map[string]string) (dataSource, error) {
	url, ok := metadata["url"]
	if !ok {
		return nil, errors.New("datasource: url required")
	}
	method, ok := metadata["method"]
	if !ok {
		return nil, errors.New("datasource: method required")
	}
	body, ok := metadata["body"]
	if !ok {
		return nil, errors.New("datasource: body required")
	}
	ds := httpDataSource{
		url:    url,
		method: method,
		body:   body,
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

// List returns a list of the available data source names.
func List() []string {
	var dataSourceNames []string
	for name := range dataSources {
		dataSourceNames = append(dataSourceNames, name)
	}
	return dataSourceNames
}
