import pytest
from app.client import GinClient
from app.models import Order, Address, Product
import httpx

@pytest.mark.asyncio
async def test_gin_client_create_order(httpx_mock):
    client = GinClient(base_url="http://test:8080")
    
    order = Order(
        order_id="TEST-001",
        customer_name="Test User",
        customer_email="test@example.com",
        address=Address(
            street="123 Test St",
            city="Test City",
            country="Test Country",
            zip_code="12345"
        ),
        products=[Product(id="PROD-1", name="Test", price=10.0, quantity=1)]
    )
    
    # Настраиваем мок-ответ для POST запроса
    httpx_mock.add_response(
        method="POST",
        url="http://test:8080/api/v1/order",
        json={"success": True, "message": "Order created successfully"},
        status_code=201
    )
    
    result = await client.create_order(order)
    assert result == {"success": True, "message": "Order created successfully"}

@pytest.mark.asyncio
async def test_gin_client_health_check(httpx_mock):
    client = GinClient(base_url="http://test:8080")
    
    # Настраиваем мок-ответ для GET запроса
    httpx_mock.add_response(
        method="GET",
        url="http://test:8080/health",
        json={"success": True, "message": "OK"},
        status_code=200
    )
    
    result = await client.health_check()
    assert result is True

@pytest.mark.asyncio
async def test_gin_client_health_check_failed(httpx_mock):
    client = GinClient(base_url="http://test:8080")
    
    # Настраиваем мок-ответ с ошибкой
    httpx_mock.add_exception(
        method="GET",
        url="http://test:8080/health",
        exception=httpx.ConnectError("Connection failed")
    )
    
    result = await client.health_check()
    assert result is False