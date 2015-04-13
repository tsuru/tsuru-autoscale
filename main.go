// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"

	"github.com/tsuru/tsuru-autoscale/alarm"
	"github.com/tsuru/tsuru-autoscale/api"
)

func main() {
	alarm.StartAutoScale()
	r := api.Router()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
