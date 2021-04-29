# End-to-end testing setup

Spin up a single node Elasticsearch cluster including Kibana to perform end-to-end testing.

Elasticsearch will be available at http://localhost:9200/ and Kibana at http://localhost:5601.

## How to

### Setup local environment
Spin up the test stack:

```
docker-compose up
```

###### In case you edit the docker-compose.yml you need to run with --build option. Like: 

```
docker-compose up --build
```

### Kishell Configuration
Once it is done, configure kishell to be able to query from the local stack.

1. Build Kishell from cmd folder with 

```
go build -v ./cmd/kishell
```

2. We need to configure the server:

```
./kishell configure --server
```

###### Put this values: 

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

Once the server is known, configure the role:

```
./kishell configure --role
```

###### Put this values: 

```
    Role name: local
    Index name: logstash-*
    Window filter time (e.g. @timestamp, modified_date): @timestamp
    Set as default? [Y/n]: 
```

Now it is possible to query the data with a known expected result thanks to the data loaded on the local stack.

```
./kishell search --newer="8760h" --query="clientip:172.155.107.128"
```

## Troubleshoot

If you have this error when you try to spin up the stack:

```
bootstrap checks failed elasticsearch | [1]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144] docker
```

Visit this link to fix it according to your OS: [Docker Client Run Mode](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-cli-run-prod-mode).
