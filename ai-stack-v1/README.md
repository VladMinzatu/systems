## Run

```
uv run fastapi dev
```
(To set that up, used `uv add fastapi --extra standard` after uv init, according to instructions [here](https://docs.astral.sh/uv/guides/integration/fastapi/#migrating-an-existing-fastapi-project)).

Testing chat endpoint:
```
curl -X POST http://localhost:8000/chat -H "Content-Type: application/json" -d '{"prompt":"Explain FastAPI in one sentence"}'
```

Testing the streaming ednpoint:
```
curl -N -X POST http://localhost:8000/chat-stream -H "Content-Type: application/json"  -d '{"prompt":"Explain FastAPI in one sentence"}'
```

### Local setup with ollama

While in production the API might probably use a vLLM inference backend, locally i'm using [ollama](https://github.com/ollama/ollama).

To set it up, install from https://ollama.com/download (tried installing it with brew and didn't fully work for me at the time).

Check installation:
```
ollama --version
```

Pull a model for testing:
```
ollama pull qwen2.5:7b
```
btw, models are stored in
```
du -sh ~/.ollama/models 
```
But they can be check and removed preferably via:
```
ollama list

 ollama rm qwen2.5:7b 
```

Anyway, then run the model:
```
ollama run qwen2.5:7b
```
This starts the interactive shell as well as the server. Test the model serving endpoint:
```
curl http://localhost:11434/v1/chat/completions \                                                                            
  -X POST -H "Content-Type: application/json" \
  -d '{
    "model":"qwen2.5:7b",
    "messages":[
      {
        "role":"user",
        "content":"What is FastAPI in a sentence?"
      }
    ], "temperature":0.7, "max_tokens":256
  }'
```

And with that working, we're now ready to start the FastAPI server and expect it to work end to end.

### Using Docker

```
docker build -t fastapi-app .
```

```
docker run -p 8000:80 fastapi-app
```