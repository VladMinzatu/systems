import torch.nn as nn
import torch

def main():
    linear = nn.Linear(10, 5)
    print("Linear layer:", linear)

    relu = nn.ReLU()
    print("ReLU activation:", relu)

    model = nn.Sequential(
        nn.Linear(10, 5),
        nn.ReLU(),
        nn.Linear(5, 2)
    )
    print("Sequential model:", model) 
    print("Model parameters:")
    for name, param in model.named_parameters():
        print(f"  {name}: {param.shape}")

    # A model can be called like a function to perform a forward pass
    input_data = torch.randn(1, 10)  # Batch size of 1
    output = model(input_data)
    print("Model output:", output)

if __name__ == "__main__":
    main()
