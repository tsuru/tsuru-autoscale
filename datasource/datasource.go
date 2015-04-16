// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package datasource

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	Register("http", httpDataSourceFactory)
}

// dataSource represents a data source.
type dataSource interface {
	// Get gets data from data source.
	Get() (string, error)
}

type dataSourceFactory func(metadata map[string]interface{}) (dataSource, error)

var dataSources = make(map[string]dataSourceFactory)

// Register registers a new dataSource.
func Register(name string, ds dataSourceFactory) {
	dataSources[name] = ds
}

type Instance struct {
	Name     string
	Metadata map[string]interface{}
}

// New creates a new data source instance.
func New(name string, metadata map[string]interface{}) (dataSource, error) {
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

func (i *Instance) Get() (string, error) {
	ds, err := dataSources["http"](i.Metadata)
	if err != nil {
		return "", err
	}
	return ds.Get()
}

func Get(name string) (*Instance, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var i Instance
	err = conn.DataSources().Find(bson.M{"name": name}).One(&i)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

type httpDataSource struct {
	url     string
	method  string
	body    string
	headers map[string]string
}

func httpDataSourceFactory(metadata map[string]interface{}) (dataSource, error) {
	url, ok := metadata["url"].(string)
	if !ok {
		return nil, errors.New("datasource: url required")
	}
	method, ok := metadata["method"].(string)
	if !ok {
		return nil, errors.New("datasource: method required")
	}
	body, ok := metadata["body"].(string)
	if !ok {
		return nil, errors.New("datasource: body required")
	}
	headers := map[string]string{}
	if _, ok := metadata["headers"].(map[string]string); ok {
		headers = metadata["headers"].(map[string]string)
	}
	ds := httpDataSource{
		url:     url,
		method:  method,
		body:    body,
		headers: headers,
	}
	return &ds, nil
}

func (ds *httpDataSource) Get() (string, error) {
	req, err := http.NewRequest(ds.method, ds.url, strings.NewReader(ds.body))
	if err != nil {
		return "", err
	}
	for key, value := range ds.headers {
		req.Header.Add(key, value)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Types returns a list of the available data source types.
func Types() []string {
	var dataSourceNames []string
	for name := range dataSources {
		dataSourceNames = append(dataSourceNames, name)
	}
	return dataSourceNames
}
