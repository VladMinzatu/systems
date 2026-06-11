import random

from pydantic_ai import Agent, RunContext

class DiceGame:
    def __init__(self, model):
        self.agent = Agent(  
            model,
            deps_type=str,
            instructions="""You are a dice rolling game.
            The user will provide a guess for the dice roll and their name is accessible via the get_player_name tool.
            You will roll the dice using the roll_dice tool, compare it to the guess and tell the user whether they guessed or not, including their name in the response.
            """
            )

        @self.agent.tool_plain
        def roll_dice() -> int:
            return random.randint(1, 6)

        @self.agent.tool
        def get_player_name(ctx: RunContext[str]) -> str:
            return ctx.deps

    def play(self, guess: int, name: str):
        result = self.agent.run_sync(f'My guess is {guess}', deps=name)  
        print(result.output)

        print(result.all_messages())
