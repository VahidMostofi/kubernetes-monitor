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

```
data_total = from(bucket: "general")
  |> range(start: 2020-12-26T18:35:00.506627749Z, stop: 2020-12-26T18:40:00.506627749Z)
  |> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
  |> filter(fn: (r) => r["_field"] == "resource_limits_millicpu_units")
  |> filter(fn: (r) => r["state"] == "running")
  |> rename(columns: {pod_name: "podName"})
  |> filter(fn: (r) => r["podName"] =~ /^auth*/)
  |> keep(columns: ["_time","_value","podName"])

data_usage = from(bucket: "general")
  |> range(start: 2020-12-26T18:31:00.506627749Z, stop: 2020-12-26T18:40:00.506627749Z)
  |> filter(fn: (r) => r["_measurement"] == "resource_usage")
  |> filter(fn: (r) => r["_field"] == "cpu")
  |> aggregateWindow(every: 10s, fn: mean)
  |> keep(columns: ["_time","_value","podName"])

joined = join(
  tables: {d1: data_total, d2: data_usage},
  on: ["_time","podName"], method: "inner"
)
  |> filter(fn: (r) => (exists r["_value_d1"]) and (exists r["_value_d2"]))
  |> map(fn:(r) => ({ r with _value_d1: float(v: r._value_d1) }))
  |> map(fn: (r) => ({ r with _value: (r._value_d2 / r._value_d1 )* 100.0 }))
  |> group(columns: ["_time", "podName"], mode: "by")
  |> group()
  |> aggregateWindow(every: 10s, fn: mean)
joined
```