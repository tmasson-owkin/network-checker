# network-checker

```log
# HELP monitoring_status current network status
# TYPE monitoring_status gauge
monitoring_status{ip="127.0.0.1",port="443",status="false"} 1
monitoring_status{ip="127.0.0.1",port="53",status="false"} 1
monitoring_status{ip="127.0.0.1",port="80",status="false"} 1
monitoring_status{ip="8.8.8.8",port="443",status="true"} 1
monitoring_status{ip="8.8.8.8",port="53",status="true"} 1
monitoring_status{ip="8.8.8.8",port="80",status="false"} 1
# HELP monitoring_up was the last network checker query successful
# TYPE monitoring_up gauge
monitoring_up
```