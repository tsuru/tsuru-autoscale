// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	return mux.NewRouter()
}
