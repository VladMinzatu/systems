from agents.models import set_up_local_model
from applications.city_picker import CityPicker

def main():
    ollama_model = set_up_local_model()
    city_picker = CityPicker(ollama_model)
    city_picker.pick()

if __name__ == "__main__":
    main()
