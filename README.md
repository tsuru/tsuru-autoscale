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

### Wizard

Wizard is an easy way to use autoscale with `tsuru`. Wizard creates the alarms
for scale up and scale down, based on simple inputs like: ``

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

## Configuring Wizard to works with tsuru

To `wizard` works fine with `tsuru` it is necessary to configure some data sources
and the actions to scale up and scale down.

### Add the scale up action

```
curl -XPOST -d '{"name": "scale_up", "url": "http://<tsuru_url>/apps/{app}/units", "method": "PUT", "body": "units={step}&process={process}", "headers": {"Authorization": "bearer <tsuru-token>", "Content-Type": "application/x-www-form-urlencoded"}}' -H "Content-Type: application/json" <autoscale-url>/alarm
```

### Add the scale down action

```
curl -XPOST -d '{"name": "scale_down", "url": "http://<tsuru_url>/apps/{app}/units?units={step}&process={process}", "method": "DELETE", "headers": {"Authorization": "bearer <tsuru-token>", "Content-Type": "application/x-www-form-urlencoded"}}' -H "Content-Type: application/json" <autoscale-url>/alarm
```

### Add data source to get the number of units

```
curl -XPOST -d '{"name": "units", "url": "http://<tsuru_url>/apps/{app}", "method": "GET", "headers" : {"Authorization": "bearer <tsuru_token>"}}' -H "Content-Type: application/json" <autoscale-url>/datasource
```

### Add data source to get cpu data from ElasticSearch

Only configure it if you are using ElasticSearch as tsuru metrics backend.

```
curl -XPOST -d '{"name": "cpu", "url": "http://<elasticsearch_url>/<elasticsearch_index>/cpu_max/_search", "method": "POST", "body" : "{\"size\":0, \"query\": {\"filtered\": {\"filter\": {\"bool\": {\"must\": [{\"range\": {\"value\": {\"lt\": 500}}},{ \"term\": {\"app.raw\": \"{app}\"}}, {\"term\": {\"process.raw\": \"{process}\"}}]}}}}, \"aggs\": {\"range\": {\"date_range\": {\"field\": \"@timestamp\", \"ranges\": [{\"from\": \"now-5m/m\", \"to\": \"now\"}]}, \"aggs\": {\"date\": {\"date_histogram\": {\"field\": \"@timestamp\", \"interval\": \"1m\"}, \"aggs\": {\"max\": {\"max\": {\"field\": \"value\"}}, \"avg\": {\"avg\": {\"field\": \"value\"}}}}}}}}", "public": true}' -H "Content-Type: application/json" <autoscale-url>/datasource
```
