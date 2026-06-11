from pydantic_ai.models.ollama import OllamaModel
from pydantic_ai.providers.ollama import OllamaProvider

def set_up_local_model():
    model = OllamaModel(
            "qwen2.5:7b",
            provider=OllamaProvider(
            base_url="http://localhost:11434/v1"
            ),
        )
    return model