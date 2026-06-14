## Flight Recorder Demo

Single-file server setup to demo the Go Flight Recorder (since Go 1.25). Run the server:
```
go run main.go
```

The server has one endpoint `/work`, which does light work on each request when invoked, but if invoked with the query parameter: `/work?type=heavy` it will process a heavier request, which should trigger detection of the slow request, at which point the server decides to write the trace capturing the last (up to) configured time and size of data.

Keep an ongoing stream of light requests going with:
```
hey -z 10m -q 10 http://127.0.0.1:8080/work
```

And trigger a heavy request:
```
curl -i "localhost:8080/work?type=heavy"
```

And that should produce the trace.out.

## Observations

When running the light requests, we get the following output:
```
Summary:
  Total:	20.1284 secs
  Slowest:	0.0346 secs
  Fastest:	0.0008 secs
  Average:	0.0157 secs
  Requests/sec:	499.2948
  

Response time histogram:
  0.001 [1]	|
  0.004 [626]	|■■■■■■■■■■■■■
  0.008 [846]	|■■■■■■■■■■■■■■■■■■
  0.011 [972]	|■■■■■■■■■■■■■■■■■■■■■
  0.014 [1495]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.018 [1844]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.021 [1861]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.024 [1640]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.028 [698]	|■■■■■■■■■■■■■■■
  0.031 [51]	|■
  0.035 [16]	|


Latency distribution:
  10% in 0.0056 secs
  25% in 0.0111 secs
  50% in 0.0164 secs
  75% in 0.0209 secs
  90% in 0.0239 secs
  95% in 0.0252 secs
  99% in 0.0272 secs
  ```

But when making the heavy request, we trigger our conditional trace dump:
```
2026/06/14 11:58:46 Slow request detected! Request took 739.60275ms, writing trace...
2026/06/14 11:58:46 Trace written successfully
```

To inspect the trace output:
```
go tool trace trace.out
```
