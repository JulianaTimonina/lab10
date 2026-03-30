from pydantic import BaseModel, Field, EmailStr
from typing import List, Optional

class Address(BaseModel):
    street: str = Field(..., description="Street address")
    city: str = Field(..., description="City name")
    country: str = Field(..., description="Country name")
    zip_code: str = Field(..., description="Postal code")

class Product(BaseModel):
    id: str = Field(..., description="Product ID")
    name: str = Field(..., description="Product name")
    price: float = Field(..., gt=0, description="Product price")
    quantity: int = Field(..., ge=0, description="Product quantity")

class Order(BaseModel):
    order_id: str = Field(..., description="Order ID")
    customer_name: str = Field(..., description="Customer full name")
    customer_email: EmailStr = Field(..., description="Customer email")
    address: Address = Field(..., description="Shipping address")
    products: List[Product] = Field(..., min_length=1, description="Order items")
    total_amount: Optional[float] = Field(None, description="Total order amount")
    status: Optional[str] = Field(None, description="Order status")

class Response(BaseModel):
    success: bool
    message: Optional[str] = None
    data: Optional[dict] = None