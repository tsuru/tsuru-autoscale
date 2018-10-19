// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"log"
	"os"

	"github.com/getsentry/raven-go"
)

func init() {
	if sentryDSN := os.Getenv("SENTRY_DSN"); sentryDSN != "" {
		raven.SetDSN(sentryDSN)
	}
}

var lg *Logger

// Logger represents a logger
type Logger struct {
	lg *log.Logger
}

// Print writes an info in the log
func (l *Logger) Print(v ...interface{}) {
	l.lg.Print(v...)
}

// Printf writes an info with format in the log
func (l *Logger) Printf(format string, v ...interface{}) {
	l.lg.Printf(format, v...)
}

// Error writes an error in the log
func (l *Logger) Error(err error) {
	raven.CaptureError(err, nil)
	l.Print(err)
}

// New returns a new Logger
func New() *Logger {
	return &Logger{lg: log.New(os.Stdout, "[autoscale] ", 0)}
}

// Log returns the Logger
func Log() *Logger {
	if lg == nil {
		lg = New()
	}
	return lg
}
