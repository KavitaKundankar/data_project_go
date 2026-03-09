import logging
import asyncio
from fastapi import FastAPI
from services.rabbitmq import start_consumers
from core.config import settings

# Setup logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(title="RabbitMQ to HTTP Microservice")

# In-memory storage for last received data (for testing visibility)
last_received_data = {
    "data1": None,
    "data2": None
}

@app.on_event("startup")
async def startup_event():
    logger.info("Starting up microservice...")
    # Background task to run consumers
    app.state.rabbitmq_connection = await start_consumers()

@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Shutting down microservice...")
    await app.state.rabbitmq_connection.close()

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "rabbit-http-bridge"}

# Mock endpoints for testing
@app.post("/api/data1")
async def mock_api1_post(data: dict):
    logger.info(f"--- DATA RECEIVED AT HTTP1 ---")
    # logger.info(f"JSON Data: {data}")
    last_received_data["data1"] = data
    return {"status": "success", "received": data}

@app.get("/api/data1")
async def mock_api1_get():
    if not last_received_data["data1"]:
        return {"message": "No data received yet for Queue 1. Push a message to RabbitMQ!"}
    return last_received_data["data1"]

@app.post("/api/data2")
async def mock_api2_post(data: dict):
    logger.info(f"--- DATA RECEIVED AT HTTP2 ---")
    # logger.info(f"JSON Data: {data}")
    last_received_data["data2"] = data
    return {"status": "success", "received": data}

@app.get("/api/data2")
async def mock_api2_get():
    if not last_received_data["data2"]:
        return {"message": "No data received yet for Queue 2. Push a message to RabbitMQ!"}
    return last_received_data["data2"]

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=settings.API_PORT)
