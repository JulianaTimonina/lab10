from jose import JWTError, jwt
from typing import Optional, Dict
import datetime

SECRET_KEY = "my-secret-key-change-in-production"
ALGORITHM = "HS256"

def validate_token(token: str) -> Optional[Dict]:
    """Validate JWT token from Go service"""
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        return payload
    except JWTError:
        return None

def get_username_from_token(token: str) -> Optional[str]:
    """Extract username from token"""
    payload = validate_token(token)
    if payload:
        return payload.get("username")
    return None