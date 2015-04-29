// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/tsuru/tsuru-autoscale/alarm"
	"github.com/tsuru/tsuru-autoscale/api"
)

func port() string {
	var p string
	if p = os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func runServer(c *cli.Context) {
	r := api.Router()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port()), nil))
}

func main() {
	alarm.StartAutoScale()
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "binding",
			Value: "0.0.0.0:8778",
			Usage: "binding address",
		},
		cli.StringFlag{
			Name: "mongodb-url",
		},
		cli.StringFlag{
			Name: "mongodb-database",
		},
		cli.StringFlag{
			Name: "mongodb-prefix",
		},
	}
	app.Version = "0.0.1"
	app.Name = "autoscale"
	app.Usage = "autoscale api"
	app.Action = runServer
	app.Run(os.Args)
}
