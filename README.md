 # acp-opensearch-alias-exporter
 
 This is a Prometheus exporter which exports the metric opensearch_alias_rate{namespace}
 
 ## Usage
 
 ``` 
./opensearch-reporter
 ```

Environment variables that need to be exported:
1. `OPENSEARCH_HOST`: The host of the Opensearch cluster including protocol schema without any trailing slashes
2. `USERNAME`: The username used for basic auth
3. `PASSWORD`: The password used for basic auth
