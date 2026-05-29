import torch.nn as nn
import torch
import os

model_file_name = "my_model_state_dict.pt"

# This demo shows how to save and load a PyTorch model using state_dict, which is the recommended approach.
def main():
    model = nn.Sequential(
        nn.Linear(10, 5),
        nn.ReLU(),
        nn.Linear(5, 2)
    )
    print("Original model:", model)

    torch.save(model.state_dict(), model_file_name)
    print("Model state_dict saved to '", model_file_name, "'")

    loaded_model = nn.Sequential(
        nn.Linear(10, 5),
        nn.ReLU(),
        nn.Linear(5, 2)
    )

    loaded_model.load_state_dict(torch.load(model_file_name))
    print("Loaded model:", loaded_model)

    loaded_model.eval()  # Set the model to evaluation mode
    input_data = torch.randn(1, 10)
    output = loaded_model(input_data)
    print("Output from loaded model:", output)
    
    os.remove(model_file_name)

if __name__ == "__main__":
    main()
