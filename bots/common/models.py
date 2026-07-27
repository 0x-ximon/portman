from datetime import datetime
from decimal import Decimal
from enum import StrEnum
from uuid import UUID

from pydantic import BaseModel, ConfigDict, Field


class UserRole(StrEnum):
    REGULAR = "REGULAR"
    AUTOMATED = "AUTOMATED"
    ADMINISTRATOR = "ADMINISTRATOR"


class User(BaseModel):
    id: UUID

    first_name: str = Field(alias="first_name")
    last_name: str = Field(alias="last_name")
    email_address: str = Field(alias="email_address")
    wallet_address: str = Field(alias="wallet_address")

    password: str = Field(alias="password")
    role: UserRole = Field(alias="role")
    api_key: str | None = Field(alias="api_key")

    created_at: datetime = Field(alias="created_at")
    updated_at: datetime = Field(alias="updated_at")
    deleted_at: datetime | None = Field(default=None, alias="deleted_at")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )


class Balance(BaseModel):
    id: int
    user_id: UUID = Field(alias="user_id")
    asset_id: int = Field(alias="asset_id")

    available: Decimal = Field(alias="available", ge=Decimal("0.0"))
    locked: Decimal = Field(alias="locked", ge=Decimal("0.0"))

    created_at: datetime = Field(alias="created_at")
    updated_at: datetime = Field(alias="updated_at")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )


class OrderSide(StrEnum):
    BUY = "BUY"
    SELL = "SELL"


class OrderMode(StrEnum):
    GTC = "GTC"
    FOK = "FOK"
    IOC = "IOC"


class OrderStatus(StrEnum):
    PENDING = "PENDING"
    REJECTED = "REJECTED"
    PARTIAL = "PARTIAL"
    CANCELLED = "CANCELLED"
    FULFILLED = "FULFILLED"


class Order(BaseModel):
    id: int
    user_id: UUID = Field(alias="user_id")
    ticker_id: int = Field(alias="ticker_id")

    price: Decimal = Field(alias="price", ge=Decimal("0.0"))
    quantity: Decimal = Field(alias="quantity", ge=Decimal("0.0"))
    side: OrderSide = Field(alias="side")
    mode: OrderMode = Field(alias="mode")
    status: OrderStatus = Field(alias="status")

    created_at: datetime = Field(alias="created_at")
    updated_at: datetime = Field(alias="updated_at")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )


class TickerStatus(StrEnum):
    OPEN = "OPEN"
    CLOSED = "CLOSED"
    SUSPENDED = "SUSPENDED"


class Ticker(BaseModel):
    id: int

    symbol: str = Field(alias="symbol")
    lot_size: Decimal = Field(alias="lot_size")
    tick_size: Decimal = Field(alias="tick_size")
    ask: Decimal = Field(alias="ask")
    bid: Decimal = Field(alias="bid")

    base_asset: int = Field(alias="base")
    quote_asset: int = Field(alias="quote")
    status: TickerStatus = Field(alias="status")

    created_at: datetime = Field(alias="created_at")
    updated_at: datetime = Field(alias="updated_at")
    deleted_at: datetime | None = Field(default=None, alias="deleted_at")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )


class AssetKind(StrEnum):
    CRYPTO = "CRYPTO"
    STOCK = "STOCK"
    FIAT = "FIAT"
    COMMODITY = "COMMODITY"


class Asset(BaseModel):
    id: int

    name: str = Field(alias="name")
    symbol: str = Field(alias="symbol")
    decimals: int = Field(alias="decimals")
    kind: AssetKind = Field(alias="kind")

    created_at: datetime = Field(alias="created_at")
    updated_at: datetime = Field(alias="updated_at")
    deleted_at: datetime | None = Field(default=None, alias="deleted_at")

    model_config = ConfigDict(
        populate_by_name=True,
        from_attributes=True,
    )
