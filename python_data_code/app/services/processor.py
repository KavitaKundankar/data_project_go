import httpx
from datetime import datetime
from models.schemas import QueueData, ProcessedData
from core.config import settings
import logging

logger = logging.getLogger(__name__)

async def function1(data: QueueData) -> dict:
    logger.info(f"Processing data from Queue 1: {data.id}")
    return data.payload

async def function2(data: QueueData) -> dict:
    logger.info(f"Processing data from Queue 2: {data.id}")
    return data.payload

async def send_to_http1(data: dict):
    async with httpx.AsyncClient() as client:
        try:
            response = await client.post(settings.HTTP_1_URL, json=data)
            response.raise_for_status()
            logger.info(f"Successfully sent data from Queue 1 to HTTP 1")
        except Exception as e:
            logger.error(f"Failed to send data to HTTP 1: {e}")
            raise

async def send_to_http2(data: dict):
    async with httpx.AsyncClient() as client:
        try:
            response = await client.post(settings.HTTP_2_URL, json=data)
            response.raise_for_status()
            logger.info(f"Successfully sent data from Queue 2 to HTTP 2")
        except Exception as e:
            logger.error(f"Failed to send data to HTTP 2: {e}")
            raise
