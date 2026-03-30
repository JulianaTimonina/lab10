import httpx
import os
from typing import Optional
from .models import Order, Response

class GinClient:
    def __init__(self, base_url: str = None):
        self.base_url = base_url or os.getenv("GIN_SERVICE_URL", "http://localhost:8080")
        self.client = httpx.AsyncClient(timeout=30.0)
    
    async def close(self):
        await self.client.aclose()
    
    async def create_order(self, order: Order) -> dict:
        """Send order to Gin service"""
        response = await self.client.post(
            f"{self.base_url}/api/v1/order",
            json=order.model_dump(exclude_none=True)
        )
        response.raise_for_status()
        return response.json()
    
    async def get_order(self, order_id: str) -> dict:
        """Get order from Gin service"""
        response = await self.client.get(
            f"{self.base_url}/api/v1/order/{order_id}"
        )
        response.raise_for_status()
        return response.json()
    
    async def list_orders(self) -> dict:
        """List all orders from Gin service"""
        response = await self.client.get(
            f"{self.base_url}/api/v1/orders"
        )
        response.raise_for_status()
        return response.json()
    
    async def health_check(self) -> bool:
        """Check if Gin service is healthy"""
        try:
            response = await self.client.get(f"{self.base_url}/health")
            return response.status_code == 200
        except:
            return False