// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "web/templates/index.html", nil)
}
