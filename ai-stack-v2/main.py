from pydantic_ai import Agent
from pydantic_ai.models.ollama import OllamaModel
from pydantic_ai.providers.ollama import OllamaProvider

def main():
    model = OllamaModel(
        "qwen2.5:7b",
        provider=OllamaProvider(
        base_url="http://localhost:11434/v1"
        ),
    )

    agent = Agent(  
        model,
        instructions='Be concise, reply with one sentence.',  
    )

    result = agent.run_sync('What is pydantic.ai?')  
    print(result.output)

if __name__ == "__main__":
    main()
