from agents.models import set_up_local_model
from applications.mcp_summer import AgenticSummer

def main():
    ollama_model = set_up_local_model()
    summer_agent = AgenticSummer(ollama_model)
    result = summer_agent.run("What is 5 + 7?")
    print(f"Result of 5 + 7: {result.output}")

if __name__ == "__main__":
    main()
