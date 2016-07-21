// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// FindServiceInstance returns an auto scale instance
func FindServiceInstance(token string) ([]Instance, error) {
	tsuruHost := os.Getenv("TSURU_HOST")
	url := fmt.Sprintf("%s/services/autoscale", tsuruHost)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger().Error(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode > 399 {
		logger().Printf("Got error find service instance status code > 399: body: %s. url: %s. status code: %d. request: %#v", string(body), url, resp.StatusCode, req)
		return nil, errors.New(string(body))
	}
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	var instances []Instance
	err = json.Unmarshal(body, &instances)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	return instances, nil
}
