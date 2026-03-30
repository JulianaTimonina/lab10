from fastapi import FastAPI, HTTPException, Depends
from fastapi.responses import JSONResponse
import logging
import uvicorn
from contextlib import asynccontextmanager
from typing import List
import httpx

from .models import Order, Response
from .client import GinClient

# Setup logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Global client instance
gin_client = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    global gin_client
    gin_client = GinClient()
    logger.info("FastAPI service started")
    yield
    # Shutdown
    if gin_client:
        await gin_client.close()
    logger.info("FastAPI service stopped")

app = FastAPI(
    title="FastAPI Order Service",
    description="Service that communicates with Gin order service",
    version="1.0.0",
    lifespan=lifespan
)

# Middleware for logging
@app.middleware("http")
async def log_requests(request, call_next):
    import time
    start_time = time.time()
    response = await call_next(request)
    process_time = time.time() - start_time
    logger.info(f"{request.method} {request.url.path} - {response.status_code} ({process_time:.3f}s)")
    return response

@app.get("/")
async def root():
    return {"message": "FastAPI Order Service", "status": "running"}

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    if gin_client:
        gin_healthy = await gin_client.health_check()
        return {
            "status": "healthy",
            "gin_service": "connected" if gin_healthy else "disconnected"
        }
    return {"status": "healthy", "gin_service": "not_initialized"}

@app.post("/api/v1/orders", response_model=Response)
async def create_order(order: Order):
    """
    Create a new order and forward to Gin service
    """
    try:
        logger.info(f"Creating order: {order.order_id}")
        result = await gin_client.create_order(order)
        return JSONResponse(content=result, status_code=201)
    except httpx.HTTPStatusError as e:
        logger.error(f"HTTP error: {e}")
        raise HTTPException(status_code=e.response.status_code, detail=e.response.text)
    except Exception as e:
        logger.error(f"Error creating order: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/v1/orders/{order_id}", response_model=Response)
async def get_order(order_id: str):
    """
    Get order by ID from Gin service
    """
    try:
        logger.info(f"Fetching order: {order_id}")
        result = await gin_client.get_order(order_id)
        return JSONResponse(content=result)
    except httpx.HTTPStatusError as e:
        if e.response.status_code == 404:
            raise HTTPException(status_code=404, detail="Order not found")
        raise HTTPException(status_code=e.response.status_code, detail=e.response.text)
    except Exception as e:
        logger.error(f"Error fetching order: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/v1/orders", response_model=Response)
async def list_orders():
    """
    List all orders from Gin service
    """
    try:
        logger.info("Fetching all orders")
        result = await gin_client.list_orders()
        return JSONResponse(content=result)
    except Exception as e:
        logger.error(f"Error fetching orders: {e}")
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000, log_level="info")