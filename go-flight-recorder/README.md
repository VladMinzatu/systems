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
