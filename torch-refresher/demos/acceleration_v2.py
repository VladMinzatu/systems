import torch.nn as nn
import torch
import timeit

def main():
    M = torch.rand((2000,2000))
    M_time = timeit.timeit(lambda: M @ M.T, number=100)
    print(f"[CPU] Time for 100 matrix multiplications: {M_time:.4f} seconds")

    if torch.backends.mps.is_available():
        M_mps = torch.rand((2000,2000), device="mps")
        M_mps_time = timeit.timeit(lambda: M_mps @ M_mps.T, number=100)
        print(f"[mps GPU] Time for 100 matrix multiplications: {M_mps_time:.4f} seconds")

    torch.mps.synchronize()  # Ensure all GPU operations are complete before exiting

if __name__ == "__main__":
    main()
