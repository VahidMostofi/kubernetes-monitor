```
kubectl proxy --port=8088 &
```
```
docker run --env-file .env --net host -v $PWD/telegraf.conf:/etc/telegraf/telegraf.conf:ro telegraf

```

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