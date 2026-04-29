import asyncio
import hashlib
import hmac
import json
import os
from decimal import Decimal
from typing import Any
from uuid import UUID

import httpx
import websockets
from pydantic import BaseModel

from common.exceptions import Code, NotFoundError, PortmanException
from common.logging import logger
from common.models import User
from common.typings import Failure, Payload, Success


class Worker:
    bot_id: int
    symbol: str
    price: Decimal

    id: UUID | None
    jwt: str | None
    user: User | None

    logger: Any
    client: httpx.AsyncClient

    def __init__(self, bot_id: int, client: httpx.AsyncClient, symbol: str):
        self.bot_id = bot_id
        self.symbol = symbol
        self.price = Decimal("0")

        self.id = None
        self.jwt = None
        self.user = None

        self.client = client
        self.logger = logger.bind(id=f"worker-{bot_id}", symbol=self.symbol)

    async def run(self):
        try:
            await self.connect()
            print(f"Bot #{self.bot_id} - Connected")

            await self.retrieve()
            print(f"Bot #{self.bot_id} - Retrieved")

            asyncio.create_task(self.handle_ticks(self.symbol))
            # await self.execute()

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

        except PortmanException:
            raise

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    async def retrieve(self) -> None:
        try:
            if self.user is not None:
                return

            headers = {"Authorization": f"Bearer {self.jwt}"}
            response = await self.client.get(f"/users/{self.id}", headers=headers)
            payload = Payload[User].model_validate(response.json())

            match payload.root:
                case Success(data=user):
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
                    self.user = user

                case Failure(error=err):
                    self.logger.error("failed to create user", code=err.code, detail=err.detail)
                    raise PortmanException.into(err.code, err.detail)

        except PortmanException:
            raise

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    async def get_user(self):
        try:
            pass

        except PortmanException:
            raise

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    # async def execute(self):
    #     if not self.user or not self.user.api_key:
    #         return Err(Exception("User not initialized"))

    #     user, api_key = self.user, self.user.api_key
    #     headers: dict[str, str] = {"X-API-KEY": api_key}

    #     while True:
    #         try:
    #             await asyncio.sleep(random.randint(0, 60))

    #             payload = self.generate_order_payload()
    #             response = await self.client.post("/orders", headers=headers, json=payload)
    #             response.raise_for_status()

    #             data = response.json()["data"]
    #             order = Order.model_validate(data)
    #             print(f"Bot #{self.id} submitted {order.side} order for {order.ticker_symbol} at {order.price}")

    #         except Exception as err:
    #             print(f"Bot #{self.id} - {type(err).__name__}: {err}")
    #             break

    async def handle_ticks(self, symbol: str) -> Result[None, Exception]:

        try:
            ws_url = str(self.client.base_url).replace("http", "ws") + "/tickers/tick"
            async with websockets.connect(ws_url) as ws:
                await ws.send("ETHUSDT")
                async for message in ws:
                    data = json.loads(message)
                    self.price = Decimal(data["last"])

            return Ok(None)

        except websockets.exceptions.ConnectionClosedError:
            print(f"Bot #{self.id} - Websocket connection closed.")
            return Err(Exception("Websocket connection closed"))

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e

    # try:

    #     response.raise_for_status()
    #     data = response.json()["data"]

    #     user = User.model_validate(data)
    #     return Ok(user)

    # except Exception as e:
    #     return Err(e)

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

    # def generate_order_payload(self) -> dict[str, Any]:
    #     assert self.user.id is not None, "user id is not set"
    #     ticker_symbol = self.symbol
    #     user_id = str(self.user.id)

    #     price = float(Decimal(self.price))
    #     quantity = float(Decimal("1.00"))

    #     side = random.choice([Side.BUY, Side.SELL])
    #     type = random.choice([Type.LIMIT, Type.MARKET])

    #     return {
    #         "user_id": user_id,
    #         "ticker_symbol": ticker_symbol,
    #         "quantity": quantity,
    #         "price": price,
    #         "side": side,
    #         "type": type,
    #     }
