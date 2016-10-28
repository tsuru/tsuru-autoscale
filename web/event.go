// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"

	"github.com/tsuru/tsuru-autoscale/alarm"
)

func eventHandler(w http.ResponseWriter, r *http.Request) error {
	events, err := alarm.FindEventsBy(nil, 1000)
	if err != nil {
		return err
	}
	return render(w, "web/templates/event/list.html", events)
}
