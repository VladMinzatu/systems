# Multi-agent example
from pydantic_ai import Agent, RunContext, UsageLimits

class CityPicker:
    def __init__(self, model):
        self.city_picker_agent = Agent(  
            model,
            instructions="""Use the `generate_city` tool to generate some cities.
            Then pick the largest one among them by population and return the name of the city.
            """
            )
        self.city_generator_agent = Agent(
            model
            )
        
        @self.city_picker_agent.tool
        def generate_city(ctx: RunContext[None], count: int) -> list[str]:
            result = self.city_generator_agent.run_sync(
                f'Generate {count} European cities and return their names as a list.',
                usage=ctx.usage
                )
            return result.output


    def pick(self):
        result = self.city_picker_agent.run_sync(
            "Pick a city",
            usage_limits=UsageLimits(request_limit=5,total_tokens_limit=1000)
            )  
        print(result.output)

        print(result.all_messages())