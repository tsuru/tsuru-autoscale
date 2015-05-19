// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func FindServiceInstance(token string) ([]Instance, error) {
	tsuruHost := os.Getenv("TSURU_HOST")
	url := fmt.Sprintf("%s/service/autoscale", tsuruHost)
	resp, err := http.Get(url)
	if err != nil {
		logger().Printf("Got error on get service instances. err: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger().Printf("Got error while parsing service json: %s", err)
		return nil, err
	}
	var instances []Instance
	err = json.Unmarshal(body, &instances)
	if err != nil {
		logger().Printf("Got error on unmarshal json %s. err: %s", string(body), err)
		return nil, err
	}
	return instances, nil
}
