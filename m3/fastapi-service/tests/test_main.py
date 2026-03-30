import pytest
from fastapi.testclient import TestClient
from unittest.mock import AsyncMock, patch
import httpx  # <-- ДОБАВЬТЕ ЭТУ СТРОКУ
from app.main import app

client = TestClient(app)

def test_root():
    response = client.get("/")
    assert response.status_code == 200
    assert response.json()["message"] == "FastAPI Order Service"

def test_health_check():
    with patch("app.main.gin_client") as mock_client:
        mock_client.health_check = AsyncMock(return_value=True)
        response = client.get("/health")
        assert response.status_code == 200
        data = response.json()
        assert data["status"] == "healthy"
        assert data["gin_service"] == "connected"

@pytest.mark.asyncio
async def test_create_order():
    order_data = {
        "order_id": "TEST-001",
        "customer_name": "Test User",
        "customer_email": "test@example.com",
        "address": {
            "street": "123 Test St",
            "city": "Test City",
            "country": "Test Country",
            "zip_code": "12345"
        },
        "products": [
            {"id": "PROD-1", "name": "Test Product", "price": 99.99, "quantity": 1}
        ]
    }
    
    with patch("app.main.gin_client") as mock_client:
        mock_client.create_order = AsyncMock(return_value={
            "success": True,
            "message": "Order created successfully",
            "data": order_data
        })
        
        response = client.post("/api/v1/orders", json=order_data)
        assert response.status_code == 201
        assert response.json()["success"] is True

def test_create_order_validation_error():
    invalid_order = {
        "order_id": "TEST-002",
        # Missing required fields
    }
    
    response = client.post("/api/v1/orders", json=invalid_order)
    assert response.status_code == 422  # Validation error

@pytest.mark.asyncio
async def test_get_order():
    order_id = "TEST-003"
    
    with patch("app.main.gin_client") as mock_client:
        mock_client.get_order = AsyncMock(return_value={
            "success": True,
            "data": {"order_id": order_id, "customer_name": "Test"}
        })
        
        response = client.get(f"/api/v1/orders/{order_id}")
        assert response.status_code == 200
        assert response.json()["success"] is True

@pytest.mark.asyncio
async def test_get_order_not_found():
    order_id = "NONEXISTENT"
    
    with patch("app.main.gin_client") as mock_client:
        # Создаем мок для HTTP ошибки
        mock_request = httpx.Request("GET", f"/orders/{order_id}")
        mock_response = httpx.Response(404, request=mock_request)
        
        mock_client.get_order = AsyncMock(side_effect=httpx.HTTPStatusError(
            "Not Found",
            request=mock_request,
            response=mock_response
        ))
        
        response = client.get(f"/api/v1/orders/{order_id}")
        assert response.status_code == 404
        assert response.json()["detail"] == "Order not found"

@pytest.mark.asyncio
async def test_list_orders():
    with patch("app.main.gin_client") as mock_client:
        mock_client.list_orders = AsyncMock(return_value={
            "success": True,
            "data": [
                {"order_id": "ORD-1", "customer_name": "User 1"},
                {"order_id": "ORD-2", "customer_name": "User 2"}
            ]
        })
        
        response = client.get("/api/v1/orders")
        assert response.status_code == 200
        assert response.json()["success"] is True
        assert len(response.json()["data"]) == 2