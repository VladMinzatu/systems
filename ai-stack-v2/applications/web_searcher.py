from pydantic_ai import Agent
from pydantic_ai.common_tools.duckduckgo import duckduckgo_search_tool

# Uses common tool for web search: duckduckgo_search_tool, which is a wrapper around DuckDuckGo's search API.
class WebSearcher:
    def __init__(self, model):
        self.agent = Agent(  
            model,
            tools=[duckduckgo_search_tool()],
            instructions="Search DuckDuckGo for the given query and return the results."
            )

    def search(self, query: str):
        result = self.agent.run_sync(query)  
        print(result.output)

        print(result.all_messages())
