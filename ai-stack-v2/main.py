from applications.dice_game import DiceGame
from agents.models import set_up_local_model

def main():
    ollama_model = set_up_local_model()
    game = DiceGame(ollama_model)
    game.play(guess=3, name="Vlad")

if __name__ == "__main__":
    main()
