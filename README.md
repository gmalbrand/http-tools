# HTTP tools
Http tool kit for go

## HTTP Mirror

An HTTP server use for debugging and testing with access Combined log printed on stdout.
Several component can be reused like MonitoredMux, AccessCombinedLogger 

### Metrics URI 
Expose promhttp metrics
Custom metrics :
- http_requests_in_flight : A gauge of requests currently being served by the wrapped handler. 
- http_requests_total : A counter for requests to the wrapped handler.
- http_request_duration_seconds : A histogram of latencies for requests.
- http_response_size_bytes : A histogram of response sizes for requests.
- http_request_size_bytes : A histogram of requests sizes.

### Dump URI
Accept all method.
Dump the request headers and body in the response.

### cpuLoad and memLoad URIs
Generate load on the server
- cpuLoad : generate cpu usage (parameter load, percentage, default 80%) for a period of time (parameter duration, in second, default 10s)
Example: 
```
 curl http://localhost:8080/cpuLoad?load=120&duration=60
```

- memLoad : generate memory usage (parameter size, in MB, default 80) for a period of time (parameter duration, in second, default 10s)
Example:
```
 curl http://localhost:8080/cpuLoad?mem=1024&duration=60
```

### proxy URI
Simple HTTP proxy. Not tested.

### Build and run locally
Build from make file
```
make dep http-mirror
```

Run on default port 8080 
```
./bin/http-mirror
```
Run on custom port 
```
HTTP_SEVER_PORT=4242 ./bin/http-mirror
```

### Use with docker
Use the latest docker image available on docker hub.
```
    docker run -d -p 8080:8080 gmalbrand/http-mirror:latest 
```

### Deploy on k8s
The given manifest creates a deployment set with 2 replicas and load balancer service
The pods have prometheus and DataSet annoations for metric collection

```
    kubectl create namespace "your-namespace"
    kubectl apply -f ./deploy/deployment.yaml -n "your-namespace"
```

## Certificate generator
Generate CA cert and keys and self signed certificates.

```
make dep certgen
```
