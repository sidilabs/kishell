[![Actions Status](https://github.com/sidilabs/kishell/workflows/build/badge.svg)](https://github.com/sidilabs/kishell/actions)
[![Code Coverage](https://codecov.io/gh/sidilabs/kishell/branch/main/graph/badge.svg)](https://codecov.io/gh/sidilabs/kishell)
[![Go Report Card](https://goreportcard.com/badge/github.com/sidilabs/kishell)](https://goreportcard.com/report/github.com/sidilabs/kishell)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/4780/badge)](https://bestpractices.coreinfrastructure.org/projects/4780)
# kishell
Data export CLI for Elasticsearch

## Build

```
go build -v github.com/sidilabs/kishell/cmd/kishell
```

## Usage

Read the manual:
```
./kishell configure -h
```
```
Usage: kishell configure

Init ES server configs

Flags:
  -h, --help      Show context-sensitive help.
      --debug     Enable debug mode.

      --server    Add a new server definition
      --role      Add a new role definition
      --reset     Reset the whole configuration
```
Add a server to the configuration:
```
./kishell configure --server
```
Example given:
```
    Server name: local
    Protocol: http
    Hostname: localhost
    Port: 5601
    Username: 
    Password: 
    Kibana Version: 6.8.6
    Set as default? [Y/n]: 
```

Define the role to be used:
```
./kishell configure --role
```
Example given:
```
    Role name: local
    Index name: logstash-*
    Window filter time (e.g. @timestamp, modified_date): @timestamp
    Set as default? [Y/n]: 
```

Send queries to Elasticsearch using the [query string syntax](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/query-dsl-query-string-query.html#query-string-syntax):
```
./kishell search <QUERY>
```
Example given:
```
./kishell search --newer="8760h" --query="clientip:172.155.107.128"
```