# Multi-agent example using programatic agent hand-off - agents coordinated by code or human in the loop
from pydantic_ai import Agent, RunContext, UsageLimits
from rich.prompt import Prompt
from pydantic import BaseModel
import json

class CityOptions(BaseModel):
  cities: list[str]

class CityPickerV2:
    def __init__(self, model):
        self.city_picker_agent = Agent(  
            model,
            output_type=str,
            instructions="When given a list of cities, pick one of them and return it as a string by calling the output tool."
            )
        self.city_generator_agent = Agent(
            model,
            output_type=CityOptions,
            instructions="When asked for a list of cities, call the final result tool with the CityOptions response with a list of city names. Do not output prose."
            )

    def pick(self):
        can_pick = False
        while not can_pick:
            options = self.city_generator_agent.run_sync(f'Give me 5 European cities. Think a bit outside the box, get creative, don\'t suggest the same biggest capitals each time.',)
            print(f'I propose choosing one of these cities {options.output.cities}')
            
            answered_properly = False
            while not answered_properly:
                answer = Prompt.ask(
                'Are you happy with the options? \\[yes/no]',
                )
                if answer == 'yes':
                    can_pick = True
                    answered_properly = True
                elif answer == 'no':
                    print('No worries, let\'s try that again')
                    answered_properly = True
                else:
                    print('Sir, c\'mon, please answer yes or no')

        print('Great, let me pick one for you')
        response = self.city_picker_agent.run_sync(f'Pick one of the cities from the list: {json.dumps(options.output.cities)}')
        city = response.output
        print(f'You\'re going to {city}!')
