## Run

```
uv run fastapi dev
```

Testing chat endpoint:
```
curl -X POST http://localhost:8000/chat -H "Content-Type: application/json" -d '{"prompt":"Explain FastAPI in one sentence"}'
```

### Using Docker

```
docker build -t fastapi-app .
```

```
docker run -p 8000:80 fastapi-app
```