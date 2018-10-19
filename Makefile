# Copyright 2016 tsuru-autoscale authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

build:
	go build -x

test:
	./go.test.bash

lint:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run -c ./.golangci.yml
