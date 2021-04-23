#!/bin/bash

until curl -s http://elasticsearch:9200/_cat/health -o /dev/null; do
    echo 'Waiting for Elasticsearch...'
    sleep 10
done

until curl -s http://kibana:5601/login -o /dev/null; do
    echo 'Waiting for Kibana...'
    sleep 10
done

echo 'Extra settings for single node cluster...'
curl -XPUT -H 'Content-Type: application/json' 'http://elasticsearch:9200/_settings' -d '
{
    "index" : {
        "number_of_replicas" : 0
    }
}'

echo 'Setup mappings...'
curl -X PUT "http://elasticsearch:9200/logstash-2020.05.18?pretty" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "log": {
      "properties": {
        "geo": {
          "properties": {
            "coordinates": {
              "type": "geo_point"
            }
          }
        }
      }
    }
  }
}
'

curl -X PUT "http://elasticsearch:9200/logstash-2020.05.19?pretty" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "log": {
      "properties": {
        "geo": {
          "properties": {
            "coordinates": {
              "type": "geo_point"
            }
          }
        }
      }
    }
  }
}
'

curl -X PUT "http://elasticsearch:9200/logstash-2020.05.20?pretty" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "log": {
      "properties": {
        "geo": {
          "properties": {
            "coordinates": {
              "type": "geo_point"
            }
          }
        }
      }
    }
  }
}
'
echo 'Setup mappings done...'

echo 'Load data...'
curl -s -H 'Content-Type: application/x-ndjson' -XPOST 'http://elasticsearch:9200/_bulk?pretty' --data-binary @/tmp/data/logs.jsonl
echo 'Data loaded...'

echo 'Create index pattern...'
curl -X POST "http://kibana:5601/api/saved_objects/index-pattern/logstash" -H 'kbn-xsrf: true' -H 'Content-Type: application/json' -d'
{
    "attributes": {
        "title": "logstash-*",
        "timeFieldName": "@timestamp"
    }
}'
echo 'Index pattern created...'