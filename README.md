# tsuru-autoscale

[![Build Status](https://travis-ci.org/tsuru/tsuru-autoscale.png?branch=master)](https://travis-ci.org/tsuru/tsuru-autoscale)
[![Go Report Card](https://goreportcard.com/badge/github.com/tsuru/tsuru-autoscale)](https://goreportcard.com/report/github.com/tsuru/tsuru-autoscale)

## roadmap

### 0.1

* api to manage auto scale configuration
* support tsuru metrics backend
* single action expression: "{metric} {operator} {value}"
* scale up / scale down
* integration by bind / unbind
* tsuru plugin

### 0.2

* http backend
* more than one scale (up / down) action
* auto scale disable action
