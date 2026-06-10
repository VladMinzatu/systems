## Local run

Everything is configured in main.py for now, so just:
```
uv run main.py
```

**Note**:
No need to load the model and start serving before, the script will do it, but it won't unload it. So afterwards, check what models are being served:
```
ollama ps
```

And then they can be stopped:
```
ollama stop qwen2.5:7b
```
