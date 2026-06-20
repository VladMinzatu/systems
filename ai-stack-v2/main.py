from agents.models import set_up_local_model
from applications.city_picker_v2 import CityPickerV2

def main():
    ollama_model = set_up_local_model()
    picker = CityPickerV2(ollama_model)
    result = picker.pick()

if __name__ == "__main__":
    main()
