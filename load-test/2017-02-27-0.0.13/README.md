# Run parameters

* apid version: 0.0.13
* apid config: [apid_config.yaml]() 
* aws machine: m4.large
* mockServer params: ./mockServer -numDeps=100 -numDevs=50000 -addDevEach=3s -upDevEach=1s -upDepEach=3s
* artillery scenario: [users.yaml]()

# Results

* Artillery report: [artillery_report_20170301_100486.json]()
* Influxdb: [telegraf.autogen.00002.00]()

![dashboard](Screen Shot 2017-03-01 at 10.21.34 AM.png)
