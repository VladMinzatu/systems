from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from sse_starlette.sse import EventSourceResponse
import httpx
import json

VLLM_URL = "http://localhost:11434" # For running with ollama locally

router = APIRouter()


class GenerateRequest(BaseModel):
    prompt: str

class GenerateResponse(BaseModel):
    text: str

@router.post("/chat", response_model=GenerateResponse)
async def generate_text(req: GenerateRequest):

    text = await generate(req.prompt)

    return GenerateResponse(
        text=text
    )

@router.post("/chat-stream")
async def chat_stream(req: GenerateRequest):

    async def event_generator():
        async for chunk in stream_generate(req.prompt):
            yield {
                "event": "token",
                "data": chunk,
            }
        yield {
            "event": "done",
            "data": "",
        }

    return EventSourceResponse(event_generator())

async def generate(prompt: str):
    payload = {
        "model": "qwen2.5:7b", # TODO: make this configurable - this is the ollama name for running locally
        "messages": [
            {
                "role": "user",
                "content": prompt
            }
        ],
        "temperature": 0.7,
        "max_tokens": 256
    }

    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{VLLM_URL}/v1/chat/completions",
            json=payload,
            timeout=120
        )

    response.raise_for_status()

    data = response.json()

    return data["choices"][0]["message"]["content"]

async def stream_generate(prompt: str):

    payload = {
        "model": "qwen2.5:7b",
        "messages": [
            {
                "role": "user",
                "content": prompt,
            }
        ],
        "temperature": 0.7,
        "max_tokens": 256,
        "stream": True,
    }

    async with httpx.AsyncClient(timeout=None) as client:

        async with client.stream(
            "POST",
            f"{VLLM_URL}/v1/chat/completions",
            json=payload,
        ) as response:

            response.raise_for_status()

            async for line in response.aiter_lines():

                if not line:
                    continue

                if not line.startswith("data:"):
                    continue

                data = line.removeprefix("data: ").strip()

                if data == "[DONE]":
                    break

                chunk = json.loads(data)

                token = (
                    chunk["choices"][0]
                    .get("delta", {})
                    .get("content")
                )

                if token:
                    yield token
