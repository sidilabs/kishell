[![Actions Status](https://github.com/sidilabs/kishell/workflows/build/badge.svg)](https://github.com/sidilabs/kishell/actions)
[![Code Coverage](https://codecov.io/gh/sidilabs/kishell/branch/main/graph/badge.svg)](https://codecov.io/gh/sidilabs/kishell)
[![Go Report Card](https://goreportcard.com/badge/github.com/sidilabs/kishell)](https://goreportcard.com/report/github.com/sidilabs/kishell)
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

Add a server to the configuration:
```
./kishell configure --server  
```

Define the role to be used:
```
./kishell configure --role
```