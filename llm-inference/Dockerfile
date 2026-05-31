FROM python:3.13

RUN pip install uv

WORKDIR /app

COPY . .

RUN uv pip install --system vllm

CMD ["python", "main.py"]
