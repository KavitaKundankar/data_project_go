from pydantic_settings import BaseSettings
from typing import Optional

class Settings(BaseSettings):
    RABBITMQ_URL: str = "amqp://user:password@localhost:5672/myvhost"
    QUEUE_1_NAME: str = "dummy_data_queue_1"
    QUEUE_2_NAME: str = "dummy_data_queue_2"
    
    HTTP_1_URL: str = "http://localhost:8000/api/data1"
    HTTP_2_URL: str = "http://localhost:8000/api/data2"
    
    API_PORT: int = 8000

    class Config:
        env_file = ".env"

settings = Settings()
