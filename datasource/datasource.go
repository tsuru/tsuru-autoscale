// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package datasource

import (
	"errors"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"strings"

	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru-autoscale/log"
	"gopkg.in/mgo.v2/bson"
)

func logger() *stdlog.Logger {
	return log.Logger()
}

// DataSource represents a data source.
type DataSource struct {
	Name    string
	URL     string
	Method  string
	Body    string
	Headers map[string]string
	Public  bool
}

// New creates a new data source instance.
func New(ds *DataSource) error {
	if ds.URL == "" {
		return errors.New("datasource: url required")
	}
	if ds.Method == "" {
		return errors.New("datasource: method required")
	}
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.DataSources().Insert(&ds)
}

// FindBy returns a list of data sources filtered by "query".
func FindBy(query bson.M) ([]DataSource, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var ds []DataSource
	err = conn.DataSources().Find(query).All(&ds)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// Get finds a data source by name.
func Get(name string) (*DataSource, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var ds DataSource
	err = conn.DataSources().Find(bson.M{"name": name}).One(&ds)
	if err != nil {
		return nil, err
	}
	return &ds, nil
}

// Remove removes a data source.
func Remove(ds *DataSource) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.DataSources().Remove(ds)
}

// Get tries to get the data from the data source.
func (ds *DataSource) Get(appName string) (string, error) {
	body := strings.Replace(ds.Body, "{app}", appName, -1)
	url := strings.Replace(ds.URL, "{app}", appName, -1)
	req, err := http.NewRequest(ds.Method, url, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	for key, value := range ds.Headers {
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
