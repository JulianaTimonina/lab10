from fastapi import FastAPI, HTTPException, Depends, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from contextlib import asynccontextmanager
import httpx
from pydantic import BaseModel
from auth.jwt_validator import validate_token

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    print("Starting up Python service...")
    yield
    # Shutdown
    print("Shutting down Python service...")

app = FastAPI(
    title="Python JWT Client", 
    description="Client for Go JWT Service",
    lifespan=lifespan
)
security = HTTPBearer()

import os
GO_SERVICE_URL = os.getenv("GO_SERVICE_URL", "http://localhost:8080")

class LoginRequest(BaseModel):
    username: str
    password: str

class LoginResponse(BaseModel):
    token: str
    message: str

class ProtectedResponse(BaseModel):
    message: str
    username: str

class PublicResponse(BaseModel):
    message: str

@app.get("/")
async def root():
    return {"message": "Python JWT Client Service"}

@app.post("/login", response_model=LoginResponse)
async def login(credentials: LoginRequest):
    """Login to Go service and get JWT token"""
    async with httpx.AsyncClient(timeout=10.0) as client:
        try:
            response = await client.post(
                f"{GO_SERVICE_URL}/login",
                json=credentials.model_dump()
            )
            if response.status_code == 200:
                data = response.json()
                return LoginResponse(
                    token=data["token"],
                    message="Successfully logged in"
                )
            else:
                raise HTTPException(
                    status_code=response.status_code,
                    detail="Invalid credentials"
                )
        except httpx.RequestError as e:
            raise HTTPException(
                status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
                detail=f"Go service is unavailable: {str(e)}"
            )

@app.get("/public", response_model=PublicResponse)
async def get_public_data():
    """Get public data from Go service"""
    async with httpx.AsyncClient(timeout=10.0) as client:
        try:
            response = await client.get(f"{GO_SERVICE_URL}/public")
            if response.status_code == 200:
                return PublicResponse(**response.json())
            else:
                raise HTTPException(
                    status_code=response.status_code,
                    detail="Failed to get public data"
                )
        except httpx.RequestError as e:
            raise HTTPException(
                status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
                detail=f"Go service is unavailable: {str(e)}"
            )

@app.get("/protected", response_model=ProtectedResponse)
async def get_protected_data(credentials: HTTPAuthorizationCredentials = Depends(security)):
    """Get protected data from Go service with JWT token"""
    token = credentials.credentials
    
    # Validate token locally
    payload = validate_token(token)
    if not payload:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid or expired token"
        )
    
    # Forward request to Go service
    async with httpx.AsyncClient(timeout=10.0) as client:
        try:
            response = await client.get(
                f"{GO_SERVICE_URL}/protected/data",
                headers={"Authorization": f"Bearer {token}"}
            )
            
            if response.status_code == 200:
                data = response.json()
                return ProtectedResponse(
                    message=data["message"],
                    username=payload.get("username", "unknown")
                )
            else:
                raise HTTPException(
                    status_code=response.status_code,
                    detail="Failed to get protected data"
                )
        except httpx.RequestError as e:
            raise HTTPException(
                status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
                detail=f"Go service is unavailable: {str(e)}"
            )