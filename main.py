from fastapi import FastAPI

app = FastAPI()

@app.get("/")
async def read_root():
    return {"msg": "Welcome to FastAPI"}

# To run this project, navigate to the repository root and run `uvicorn main:app --host 0.0.0.0 --port 8000`