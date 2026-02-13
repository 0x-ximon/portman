from enum import Enum
from uuid import UUID
from datetime import datetime
from decimal import Decimal
from typing import Optional
from pydantic import BaseModel, ConfigDict, Field


class Role(str, Enum):
    REGULAR = "REGULAR"
    AUTOMATED = "AUTOMATED"
    ADMINISTRATOR = "ADMINISTRATOR"


class User(BaseModel):
    id: UUID

    first_name: str = Field(alias="first_name")
    last_name: str = Field(alias="last_name")
    phone_number: str = Field(alias="phone_number")

    email_address: str = Field(alias="email_address")
    wallet_address: str = Field(alias="wallet_address")

    free_balance: Decimal = Field(alias="free_balance")
    frozen_balance: Decimal = Field(alias="frozen_balance")

    password: str
    role: Role = Field(alias="role")
    api_key: Optional[str] = Field(alias="api_key")

    created_at: datetime = Field(alias="created_at")
    updated_at: datetime = Field(alias="updated_at")
    deleted_at: Optional[datetime] = Field(default=None, alias="deleted_at")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )


class Side(str, Enum):
    BUY = "BUY"
    SELL = "SELL"


class Type(str, Enum):
    LIMIT = "LIMIT"
    MARKET = "MARKET"


class Status(str, Enum):
    PENDING = "PENDING"
    CANCELLED = "CANCELLED"
    FULFILLED = "FULFILLED"
    REJECTED = "REJECTED"


class Order(BaseModel):
    id: int
    user_id: UUID
    ticker_symbol: str = Field(alias="ticker_symbol", min_length=3, max_length=10)

    price: Decimal = Field(alias="price", ge=Decimal("0.0"))
    quantity: Decimal = Field(alias="quantity", ge=Decimal("0.0"))

    side: Side = Field(alias="side")
    type: Type = Field(alias="type")
    status: Status = Field(alias="status")

    created_at: datetime = Field(alias="created_at")
    updated_at: datetime = Field(alias="updated_at")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )


class TickerStatus(str, Enum):
    OPEN = "OPEN"
    CLOSED = "CLOSED"
    SUSPENDED = "SUSPENDED"


class Ticker(BaseModel):
    id: int

    base: str = Field(alias="base", min_length=3, max_length=10)
    quote: str = Field(alias="quote", min_length=3, max_length=10)
    symbol: str = Field(alias="symbol", min_length=3, max_length=10)

    ask: Decimal = Field(alias="ask", ge=Decimal("0.0"))
    bid: Decimal = Field(alias="bid", ge=Decimal("0.0"))
    last: Decimal = Field(alias="last", ge=Decimal("0.0"))
    status: TickerStatus = Field(alias="status")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )
