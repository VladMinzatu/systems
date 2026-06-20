# Multi-agent example using programatic agent hand-off - agents coordinated by code or human in the loop
from pydantic_ai import Agent, RunContext, UsageLimits
from rich.prompt import Prompt
from pydantic import BaseModel

class CityOptions(BaseModel):
  cities: list[str]

class CityPickerV2:
    def __init__(self, model):
        self.city_picker_agent = Agent(  
            model,
            output_type=str,
            instructions="When given a list of cities, pick one of them and return it as a string."
            )
        self.city_generator_agent = Agent(
            model,
            output_type=CityOptions,
            instructions="""
                When asked for a list of cities, generate as many as requested and return them in a list.
                When you have the answer, call the final result tool with the CityOptions response. Do not output prose.
            """
            )

    def pick(self):
        options = self.city_generator_agent.run_sync(f'Generate 5 European cities and return their names as a list.',)
        print(options.output)

        prompt = Prompt.ask(
          'Are you happy with the options? \\[yes/no]',
        )
        print(f'You said {prompt}')
