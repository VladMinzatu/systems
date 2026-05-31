from vllm import LLM, SamplingParams

model_name = "Qwen/Qwen2.5-0.5B"

def main():
    llm = LLM(model=model_name, dtype="float16")
    params = SamplingParams(temperature=0.7, top_p=0.9, max_output_tokens=512)

    response = llm.generate("What is the capital of France?", sampling_params=params)
    for r in response:
        print(f"Generated text: {r.text}")

if __name__ == "__main__":
    main()
