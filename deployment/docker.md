## How to deploy using docker?

#### Start the proxy server
```
kubectl proxy --port=8088 &
```

#### Start the docker container with the current telegraf.conf file
```
docker run --rm --env-file .env --net host -v $PWD/telegraf.conf:/etc/telegraf/telegraf.conf:ro -v $PWD/report-metrics.out:/report-metrics.out -v $PWD/kube-configs:/root/.kube telegraf
```
