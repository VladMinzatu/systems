import timeit
from torch import nn
import torch

def main():
    model = nn.Sequential(
        nn.Linear(10, 5),
        nn.ReLU(),
        nn.Linear(5, 2)
    )
    input_data = torch.randn(1, 10)
    # Measure the time taken for a forward pass on CPU (small model)
    cpu_time = timeit.timeit(lambda: model(input_data), number=1000)
    print(f"[CPU] Small model: Time for 1000 forward passes: {cpu_time:.4f} seconds")

    # If an MPS GPU is available, measure small model there too
    if torch.backends.mps.is_available():
        model_mps = model.to("mps")
        input_mps = input_data.to("mps")
        gpu_time = timeit.timeit(lambda: model_mps(input_mps), number=1000)
        print(f"[mps GPU] Small model: Time for 1000 forward passes: {gpu_time:.4f} seconds")

    # Demo with a larger workload to better showcase GPU acceleration
    large_model = nn.Sequential(
        nn.Linear(1000, 2000),
        nn.ReLU(),
        nn.Linear(2000, 1000),
        nn.ReLU(),
        nn.Linear(1000, 512)
    )
    large_input = torch.randn(32, 1000)  # batch size 32

    large_cpu_time = timeit.timeit(lambda: large_model(large_input), number=100)
    print(f"[CPU] Large model: Time for 100 forward passes: {large_cpu_time:.4f} seconds")

    if torch.backends.mps.is_available():
        large_model_mps = large_model.to("mps")
        large_input_mps = large_input.to("mps")
        # Warm-up on device
        for _ in range(10):
            _ = large_model_mps(large_input_mps)
        large_gpu_time = timeit.timeit(lambda: large_model_mps(large_input_mps), number=100)
        print(f"[mps GPU] Large model: Time for 100 forward passes: {large_gpu_time:.4f} seconds")

if __name__ == "__main__":
    main()
