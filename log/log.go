// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"log"
	"os"
)

var lg *log.Logger

func Logger() *log.Logger {
	if lg == nil {
		lg = log.New(os.Stdout, "[alarm] ", 0)
	}
	return lg
}
