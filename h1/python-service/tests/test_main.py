import sys
import os
import pytest

# Добавляем корневую директорию проекта в путь
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from fastapi.testclient import TestClient
from main import app

client = TestClient(app)

def test_root():
    """Test root endpoint"""
    response = client.get("/")
    assert response.status_code == 200
    assert response.json() == {"message": "Python JWT Client Service"}

def test_login_valid_credentials():
    """Test login with valid credentials (requires Go service running)"""
    response = client.post(
        "/login",
        json={"username": "admin", "password": "password123"}
    )
    if response.status_code == 503:
        pytest.skip("Go service not available")
    assert response.status_code == 200
    assert "token" in response.json()
    assert response.json()["message"] == "Successfully logged in"

def test_login_invalid_credentials():
    """Test login with invalid credentials"""
    response = client.post(
        "/login",
        json={"username": "wrong", "password": "wrong"}
    )
    if response.status_code == 503:
        pytest.skip("Go service not available")
    assert response.status_code == 401

def test_public_data():
    """Test public data endpoint"""
    response = client.get("/public")
    if response.status_code == 503:
        pytest.skip("Go service not available")
    assert response.status_code == 200
    assert "message" in response.json()

def test_protected_data_without_token():
    """Test protected endpoint without token"""
    response = client.get("/protected")
    assert response.status_code == 403

def test_protected_data_with_invalid_token():
    """Test protected endpoint with invalid token"""
    response = client.get(
        "/protected",
        headers={"Authorization": "Bearer invalid.token.here"}
    )
    assert response.status_code == 401