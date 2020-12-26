```
kubectl proxy --port=8088 &
```
```
docker run --rm --env-file .env --net host -v $PWD/telegraf.conf:/etc/telegraf/telegraf.conf:ro -v $PWD/report-metrics.out:/report-metrics.out -v $PWD/kube-configs:/root/.kube telegraf
```

Generate emtpy telegraf.conf file
```
docker run --rm telegraf telegraf config > telegraf.conf
```

contents of ```kube-configs``` directory:
- ca.crt
- client.crt
- client.key
- config -> template for this is available

### Get Pods of a service:
```
from(bucket: "general")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
  |> filter(fn: (r) => r["pod_name"] =~ /auth-*/ and r["state"] == "running")
  |> keep(columns: ["pod_name", "_time"])
  |> unique(column: "pod_name")
  |> keep(columns: ["pod_name"])
```
### Get total allocated resources
```
from(bucket: "general")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
  |> filter(fn: (r) => r["_field"] == "resource_limits_millicpu_units")
  |> filter(fn: (r) => r["pod_name"] =~ /^auth-*/)
  |> group(columns: ["_time"], mode: "by")
  |> aggregateWindow(every: 5s, fn: sum, createEmpty: false)
  |> group()
  |> rename(columns: {_value: "Auth-service"})
```
To use this you also need to:
- Go into Transform section (next to query section)
- Add new Filter By Name
- Unselect stuff you don't need
- Select stuff you need

The advantage of this approach in comparision with the other approach is that here we can deploy it without needing to have a deamonset .
