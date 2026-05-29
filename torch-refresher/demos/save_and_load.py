import torch.nn as nn
import torch
import os

model_file_name = "my_model.pt"

def main():
    # Create a simple model
    model = nn.Sequential(
        nn.Linear(10, 5),
        nn.ReLU(),
        nn.Linear(5, 2)
    )
    print("Original model:", model)

    # Save the entire model to a file 
    torch.save(model, model_file_name)
    print("Model saved to '", model_file_name, "'")

    # Load the model from the saved file
    loaded_model = torch.load(model_file_name, weights_only=False)
    print("Loaded model:", loaded_model)

    loaded_model.eval()  # Set the model to evaluation mode
    input_data = torch.randn(1, 10)
    output = loaded_model(input_data)
    print("Output from loaded model:", output)
    
    os.remove(model_file_name)

if __name__ == "__main__":
    main()
