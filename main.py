from fastapi import FastAPI

app = FastAPI()

@app.get("/")
async def read_root():
    return {"msg": "Welcome to FastAPI"}
