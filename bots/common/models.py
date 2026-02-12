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
    role: Role = Field(default=Role.REGULAR)
    api_key: Optional[str] = Field(default=None, alias="api_key")

    created_at: datetime = Field(alias="created_at")
    updated_at: datetime = Field(alias="updated_at")
    deleted_at: Optional[datetime] = Field(default=None, alias="deleted_at")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )


class Status(str, Enum):
    OPEN = "OPEN"
    CLOSED = "CLOSED"
    SUSPENDED = "SUSPENDED"


class Ticker(BaseModel):
    id: int

    base: str = Field(min_length=3, max_length=10)
    quote: str = Field(min_length=3, max_length=10)
    symbol: str = Field(min_length=3, max_length=10)

    ask: Decimal = Field(default=Decimal("0.0"), ge=Decimal("0.0"))
    bid: Decimal = Field(default=Decimal("0.0"), ge=Decimal("0.0"))
    last: Decimal = Field(default=Decimal("0.0"), ge=Decimal("0.0"))
    status: Status = Field(default=Status.CLOSED)

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )
