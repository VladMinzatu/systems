import torch

def main():
    print("PyTorch version:", torch.__version__)

    print("CUDA available:", torch.cuda.is_available())
    print("MPS available:", torch.backends.mps.is_available())

    M = torch.tensor([1, 2, 3])
    print("M device:", M.device)

    M = M.to("mps")
    print("M device after to('mps'):", M.device)

    N = torch.tensor([4, 5, 6], device="mps")
    print("N device:", N.device)


if __name__ == "__main__":
    main()
