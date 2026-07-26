import asyncio
import hashlib
import hmac
import json
import os
import random
from decimal import Decimal
from typing import Any
from uuid import UUID

import httpx
import websockets
from pydantic import BaseModel

from common.exceptions import ApiError, Code, NotFoundError, PortmanException
from common.logging import logger
from common.models import Order, OrderMode, OrderSide, User
from common.typings import Failure, Payload, Success


class Worker:
    bot_id: int
    ticker_id: int
    ticker_symbol: str
    price: Decimal

    id: UUID | None
    jwt: str | None
    user: User | None

    logger: Any
    client: httpx.AsyncClient

    def __init__(self, bot_id: int, client: httpx.AsyncClient, ticker_id: int, ticker_symbol: str):
        self.bot_id = bot_id
        self.ticker_id = ticker_id
        self.ticker_symbol = ticker_symbol
        self.price = Decimal(100)

        self.id = None
        self.jwt = None
        self.user = None

        self.client = client
        self.logger = logger.bind(id=f"worker-{bot_id}", ticker_symbol=self.ticker_symbol)

    async def run(self):
        try:
            await self.connect()
            print(f"Bot #{self.bot_id} - Connected")

            asyncio.create_task(self.handle_ticks(self.ticker_symbol))
            await self.execute()

        except PortmanException:
            raise

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    async def connect(self, retry: bool = True) -> None:
        try:
            api_key = self.get_api_key()
            headers = {"X-API-KEY": api_key}

            class Creds(BaseModel):
                id: UUID
                jwt: str

            response = await self.client.post("/auth/exchange", headers=headers)
            payload = Payload[Creds].model_validate(response.json())

            match payload.root:
                case Success(data=creds):
                    if isinstance(creds, Creds):
                        self.id = creds.id
                        self.jwt = creds.jwt

                case Failure(error=err):
                    match err.code:
                        case Code.NOT_FOUND:
                            if retry:
                                await self.create_user()
                                await self.connect(retry=False)

                            else:
                                self.logger.error("authentication failed", code=err.code, detail=err.detail)
                                raise PortmanException.into(err.code, err.detail)

                        case _:
                            self.logger.error("authentication failed", code=err.code, detail=err.detail)
                            raise PortmanException.into(err.code, err.detail)

            if self.user is None:
                await self.get_user()

        except PortmanException:
            raise

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    async def get_user(self) -> None:
        try:
            if self.user is not None:
                return

            headers = {"Authorization": f"Bearer {self.jwt}"}
            response = await self.client.get(f"/users/{self.id}", headers=headers)
            payload = Payload[User].model_validate(response.json())

            match payload.root:
                case Success(data=user):
                    if isinstance(user, User):
                        self.user = user

                case Failure(error=err):
                    self.logger.error("retrieval failed", code=err.code, detail=err.detail)
                    raise PortmanException.into(err.code, err.detail)

        except PortmanException:
            raise

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    async def create_user(self):
        try:
            api_key = self.get_api_key()
            data = {
                "first_name": "Portman",
                "last_name": f"Bot_{self.bot_id}",
                "email_address": f"portman_bot#{self.bot_id}@system.local",
                "wallet_address": f"0x0000{self.bot_id}",
                "password": f"bot#{self.bot_id}_password",
                "role": "AUTOMATED",
                "api_key": api_key,
            }

            response = await self.client.post("/users", json=data)
            payload = Payload[User].model_validate(response.json())

            match payload.root:
                case Success(data=user):
                    if isinstance(user, User):
                        self.user = user

                case Failure(error=err):
                    self.logger.error("failed to create user", code=err.code, detail=err.detail)
                    raise PortmanException.into(err.code, err.detail)

        except PortmanException:
            raise

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    async def execute(self):
        try:
            assert self.id is not None, "id is not set"
            assert self.jwt is not None, "jwt is not set"

            self.logger.info("worker execution started")
            while True:
                headers = {"Authorization": f"Bearer {self.jwt}"}
                payload = self.generate_order_payload()

                response = await self.client.post("/orders", json=payload, headers=headers)
                payload = Payload[Order].model_validate(response.json())

                match payload.root:
                    case Success(data=order):
                        if isinstance(order, Order):
                            self.logger.info("submitted order", side=order.side, mode=order.mode)

                    case Failure(error=err):
                        self.logger.error("failed to create order", code=err.code, detail=err.detail)
                        raise PortmanException.into(err.code, err.detail)

                await asyncio.sleep(random.uniform(2, 5))

        except PortmanException:
            raise

        except asyncio.CancelledError:
            self.logger.info("worker execution cancelled")
            raise

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    async def handle_ticks(self, symbol: str) -> None:
        try:
            ws_url = str(self.client.base_url).replace("http", "ws") + "/tickers/tick"
            async with websockets.connect(ws_url) as ws:
                await ws.send(symbol)

                async for message in ws:
                    data = json.loads(message)
                    self.price = Decimal(data["last"])
                    self.logger.info("received", detail=str(data))

        except websockets.exceptions.ConnectionClosedError as e:
            self.logger.error("websocket connection closed", detail=str(e))
            raise ApiError("websocket connection closed") from e

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    def get_api_key(self) -> str:
        try:
            email = f"portman_bot#{self.bot_id}@system.local"
            secret = os.getenv("SYSTEM_SECRET")

            if secret is None:
                raise NotFoundError("SYSTEM_SECRET not set")

            email_bytes = email.encode("utf-8")
            secret_bytes = secret.encode("utf-8")

            h = hmac.new(secret_bytes, email_bytes, digestmod=hashlib.sha256)
            return h.hexdigest()

        except PortmanException:
            raise

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    def generate_order_payload(self) -> dict[str, Any]:
        assert self.user is not None, "user is not set"
        ticker_id = self.ticker_id
        user_id = str(self.user.id)

        price = float(Decimal(self.price))
        quantity = float(Decimal("1.00"))

        side = random.choice([OrderSide.BUY, OrderSide.SELL])
        mode = random.choice([OrderMode.GTC, OrderMode.FOK, OrderMode.IOC])

        return {
            "user_id": user_id,
            "ticker_id": ticker_id,
            "quantity": quantity,
            "price": price,
            "side": side,
            "mode": mode,
        }
