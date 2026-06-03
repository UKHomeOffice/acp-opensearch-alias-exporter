 # acp-opensearch-alias-exporter
 
This is a Prometheus exporter which exports the metric opensearch_alias_rate{namespace}.

What this does is:
1. Gets list of Aliases from Opensearch
2. Gets the [stats](https://docs.opensearch.org/2.19/api-reference/index-apis/stats/) and [latest ISM action status](https://docs.opensearch.org/2.19/im-plugin/ism/api/#explain-index) of the write index of each Alias
3. Updates the Prometheus "opensearch_alias_rollover_health" metric for each Alias using actions.failed from the ISM endpoint for that Alias: 1 = rollover attempt failure, 0 = no failure. 
4. Waits a minute
5. Compares the current total documents of the write index with the previous total documents.
6. Output this as the Prometheus "opensearch_alias_rate" metric.

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
