 # acp-opensearch-alias-exporter
 
This is a Prometheus exporter which exports the metric opensearch_alias_rate{namespace}.

What this does is:
1. Get list of Aliases from Opensearch
2. Gets the stats of the write index of each Alias
3. Waits a minute
4. Compares the stats of the current write index with it before
5. Output this as a Prometheus metrics
6. it goes into Opensearch gets the list of Aliases and their write index, it stores the total documents

## NOTE:
If the namespace is removed from the Kubernetes cluster the metric will still be outputed as there isn't any check to see whether an alias has been removed.

 ## Usage
 
 ``` 
./opensearch-reporter
 ```

Environment variables that need to be exported:
1. `OPENSEARCH_HOST`: The host of the Opensearch cluster including protocol schema without any trailing slashes
2. `USERNAME`: The username used for basic auth
3. `PASSWORD`: The password used for basic auth
