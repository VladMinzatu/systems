from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
import httpx

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
