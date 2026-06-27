### Run
```
go run .
```

### Simple test
```
curl -i -H "Content-type: application/json" -X POST -d '{"Kind":"cpu", "Size": 100}' localhost:8080/task 
```