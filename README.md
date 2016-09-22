# tsuru-autoscale

[![Build Status](https://travis-ci.org/tsuru/tsuru-autoscale.png?branch=master)](https://travis-ci.org/tsuru/tsuru-autoscale)
[![Go Report Card](https://goreportcard.com/badge/github.com/tsuru/tsuru-autoscale)](https://goreportcard.com/report/github.com/tsuru/tsuru-autoscale)
[![codecov](https://codecov.io/gh/tsuru/tsuru-autoscale/branch/master/graph/badge.svg)](https://codecov.io/gh/tsuru/tsuru-autoscale)

## features

* HTTP api to manage auto scale configuration
* integration with [tsuru](https://tsuru.io)
* Wizard that makes easy configure an auto scale for [tsuru](https://tsuru.io)
applications
* Support HTTP based data sources
* Support HTTP based actions
* Support alarm scripts in javascript

## Install as tsuru application

### Create tsuru app using Go platform

```
tsuru app-create autoscale go
```

### Configuring MongoDB

We should use environment variables to configure the database:

```
tsuru env-set "MONGODB_URL=mongodb://172.17.0.1:27017/tsuru_autoscale" -a autoscale
```

### Deploy the applications

```
tsuru app-deploy . -a autoscale
```
