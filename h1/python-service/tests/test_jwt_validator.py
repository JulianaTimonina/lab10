import sys
import os
import pytest

# Добавляем корневую директорию проекта в путь
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from auth.jwt_validator import validate_token, get_username_from_token
from jose import jwt
import time

SECRET_KEY = "my-secret-key-change-in-production"

def test_validate_token_valid():
    """Test valid token validation"""
    token = jwt.encode(
        {"username": "testuser", "exp": time.time() + 3600},
        SECRET_KEY,
        algorithm="HS256"
    )
    payload = validate_token(token)
    assert payload is not None
    assert payload["username"] == "testuser"

def test_validate_token_expired():
    """Test expired token validation"""
    token = jwt.encode(
        {"username": "testuser", "exp": time.time() - 3600},
        SECRET_KEY,
        algorithm="HS256"
    )
    payload = validate_token(token)
    assert payload is None

def test_validate_token_invalid():
    """Test invalid token validation"""
    payload = validate_token("invalid.token.here")
    assert payload is None

def test_get_username_from_token():
    """Test extracting username from token"""
    token = jwt.encode(
        {"username": "testuser", "exp": time.time() + 3600},
        SECRET_KEY,
        algorithm="HS256"
    )
    username = get_username_from_token(token)
    assert username == "testuser"

def test_get_username_from_invalid_token():
    """Test extracting username from invalid token"""
    username = get_username_from_token("invalid")
    assert username is None