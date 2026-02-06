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

    role: Role
    password: str
    api_key: Optional[str] = Field(default=None, alias="api_key")

    created_at: datetime = Field(alias="created_at")
    updated_at: datetime = Field(alias="updated_at")
    deleted_at: Optional[datetime] = Field(default=None, alias="deleted_at")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )
