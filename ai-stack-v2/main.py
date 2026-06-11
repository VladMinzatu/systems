from agents.models import set_up_local_model
from applications.web_searcher import WebSearcher

def main():
    ollama_model = set_up_local_model()
    game = WebSearcher(ollama_model)
    game.search(query="What is the capital of France?")

if __name__ == "__main__":
    main()
