# Just sums 2 numbers, lol
# Using mcp server at mcp-servers/sum-server.py - can be run independently with "uv run mcp-servers/sum-server.py stdio", but no need to do that since pydantic_ai will run it when we run this file.
from pydantic_ai import Agent
from pydantic_ai.mcp import MCPServerStdio

class AgenticSummer():
    def __init__(self, model):
        self.server = MCPServerStdio('uv', args=['run', './mcp-servers/sum-server.py', 'stdio'], timeout=10)
        self.agent = Agent(model, 
        output_type=int, # TODO: this fails half the time because the server returns a sentence - understand what pydantic does and where prompting support is needed + retry behavior.
        toolsets=[self.server])

    def run(self, query: str) -> int:
        return self.agent.run_sync(query)
