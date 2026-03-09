from pydantic import BaseModel
from typing import Dict, Any

class QueueData(BaseModel):
    id: str
    payload: Dict[str, Any]
    content_type: str = "json"

class ProcessedData(BaseModel):
    original_id: str
    processed_at: str
    data: Dict[str, Any]
