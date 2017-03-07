# Run parameters

* apid version: 0.0.13
* aws machine: m4.large
* mockServer params: `-numDeps=100 -numDevs=50000 -addDevEach=3s -upDevEach=1s -upDepEach=3s`
* artillery scenario: [users.yaml](users.yaml)

apid config:
 ```
api_listen: :9000    # note: leave api open for connections from anywhere
api_expvar_path: /expvar
events_buffer_size: 5
apigeesync_proxy_server_base: http://localhost:9001
apigeesync_snapshot_server_base: http://localhost:9001
apigeesync_change_server_base: http://localhost:9001
apidanalytics_uap_server_base: http://localhost:9001
apigeesync_consumer_key: key
apigeesync_consumer_secret: secret
apigeesync_cluster_id: cluster
log_level: info
data_trace_log_level: info
data_source: file:%s?_busy_timeout=20000
local_storage_path: /demo/data
```

# Results

* Artillery report: [artillery_report_20170301_100486.json](artillery_report_20170301_100486.json)
* Influxdb: [telegraf.autogen.00002.00](telegraf.autogen.00002.00)

### Summary
```
    "scenariosCreated": 41016,
    "scenariosCompleted": 40836,
    "requestsCompleted": 81736,
    "latency": {
      "min": 85.5,
      "max": 9807.9,
      "median": 529.1,
      "p95": 4245.1,
      "p99": 6622.4
    },
    "rps": {
      "count": 81916,
      "mean": 285.83
    },
    "scenarioDuration": {
      "min": 201.2,
      "max": 14077.2,
      "median": 1094.1,
      "p95": 7674.2,
      "p99": 9399.1
    },
    "scenarioCounts": {
      "Verified user": 36860,
      "Unverified user": 4156
    },
    "errors": {
      "ESOCKETTIMEDOUT": 24,
      "ETIMEDOUT": 180
    },
    "codes": {
      "200": 81736
    },
```

### Influx Dashboard
![dashboard](Screen Shot 2017-03-01 at 10.21.34 AM.png)
