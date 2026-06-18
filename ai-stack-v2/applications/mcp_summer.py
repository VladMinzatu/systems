# Just sums 2 numbers, lol
# Using mcp server at mcp-servers/sum-server.py - can be run independently with "uv run mcp-servers/sum-server.py stdio", but no need to do that since pydantic_ai will run it when we run this file.
from pydantic_ai import Agent
from pydantic_ai.mcp import MCPServerStdio

class AgenticSummer():
    def __init__(self, model):
        self.server = MCPServerStdio('uv', args=['run', './mcp-servers/sum-server.py', 'stdio'], timeout=10)
        self.agent = Agent(model, 
            output_type=int,
            instructions="""
            Use the `add` tool to add two numbers together.

            If you know the answer, call the final result tool with the integer.
            Do not output prose.""", # Nudge the model towards the output tool - this already helps reliability a lot from what I've seen 
            toolsets=[self.server])

    def run(self, query: str) -> int:
        return self.agent.run_sync(query)
