import asyncio
import json
import logging
import aio_pika
from core.config import settings
from models.schemas import QueueData
from services.processor import function1, function2, send_to_http1, send_to_http2

logger = logging.getLogger(__name__)

async def process_queue1_message(message: aio_pika.IncomingMessage):
    async with message.process():
        try:
            body = json.loads(message.body.decode())
            data = QueueData(id=str(message.message_id or "unknown"), payload=body)
            processed_data = await function1(data)
            await send_to_http1(processed_data)
            logger.info("Message from Queue 1 processed and sent.")
        except Exception as e:
            logger.error(f"Error processing message from Queue 1: {e}")

async def process_queue2_message(message: aio_pika.IncomingMessage):
    async with message.process():
        try:
            body = json.loads(message.body.decode())
            data = QueueData(id=str(message.message_id or "unknown"), payload=body)
            processed_data = await function2(data)
            await send_to_http2(processed_data)
            logger.info("Message from Queue 2 processed and sent.")
        except Exception as e:
            logger.error(f"Error processing message from Queue 2: {e}")

async def start_consumers():
    connection = await aio_pika.connect_robust(settings.RABBITMQ_URL)
    channel = await connection.channel()
    
    # Declare queues
    queue1 = await channel.declare_queue(settings.QUEUE_1_NAME, durable=True)
    queue2 = await channel.declare_queue(settings.QUEUE_2_NAME, durable=True)
    
    # Start consuming
    await queue1.consume(process_queue1_message)
    await queue2.consume(process_queue2_message)
    
    logger.info("RabbitMQ Consumers started...")
    return connection
