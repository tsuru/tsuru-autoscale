# Copyright 2016 tsuru-autoscale authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

test:
	go clean ./...
	go test ./...

build:
	go build -x
