from fastapi import FastAPI
from app.routers import items

app = FastAPI()
app.include_router(items.router)

@app.get("/health")
async def root():
    return {"OK": "Service healthy"}
