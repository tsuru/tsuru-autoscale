# tsuru-autoscale

[![Build Status](https://travis-ci.org/tsuru/tsuru-autoscale.png?branch=master)](https://travis-ci.org/tsuru/tsuru-autoscale)
[![Go Report Card](https://goreportcard.com/badge/github.com/tsuru/tsuru-autoscale)](https://goreportcard.com/report/github.com/tsuru/tsuru-autoscale)
[![codecov](https://codecov.io/gh/tsuru/tsuru-autoscale/branch/master/graph/badge.svg)](https://codecov.io/gh/tsuru/tsuru-autoscale)

## Features

* HTTP api to manage auto scale configuration
* integration with [tsuru](https://tsuru.io)
* Wizard that makes easy configure an auto scale for [tsuru](https://tsuru.io)
applications
* Support HTTP based data sources
* Support HTTP based actions
* Support alarm scripts in javascript

## Concepts

The `tsuru-autoscale` is based on three elements: `data sources`, `actions` and `alarms`.

### Data sources

Data source is a http endpoint that provide the data to an alarm. Is based on the
data source data that the alarm will execute an action.

### Actions

Action is a http endpoint that is called when the alarm expression result is `true`.

### Alarms

Alarm is composed by data sources, actions and by an expression. When the expression result is `true` the actions will be executed.

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

## API Reference

### list data sources

```
curl <autoscale-url>/datasource
```

### add data source

```
curl -XPOST -d '{}' -H "Content-Type: application/json" <autoscale-url>/datasource
```

### remove a data source

```
curl -XDELETE <autoscale-url>/datasource/{name}
```

### list actions

```
curl <autoscale-url>/action
```

### add an action

```
curl -XPOST -d '{}' -H "Content-Type: applicaiton/json" <autoscale-url>/action
```

### remove an action

```
curl -XDELETE <autoscale-url>/action/{name}
```

### list alarms

```
curl <autoscale-url>/alarm
```

### add an alarm

```
curl -XPOST -d '{}' -H "Content-Type: application/json" <autoscale-url>/alarm
```

### remove an alarm

```
curl -XDELETE <autoscale-url>/alarm/{name}
```
