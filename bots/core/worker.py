import hashlib
import hmac
import httpx
import os

from common.typings import Result, Ok, Err
from common.models import User


class Worker:
    id: int
    user: User
    client: httpx.AsyncClient

    def __init__(self, bot_id: int, client: httpx.AsyncClient):
        self.id = bot_id
        self.client = client

    async def run(self):
        result = await self.connect()
        match result:
            case Ok(_):
                print(f"Bot #{self.id} - Connected")
            case Err(error):
                print(f"Bot #{self.id} - Failed to connect: {error}")

    async def connect(self) -> Result[None, Exception]:
        result = await self.get_user()
        match result:
            case Ok(data):
                self.user = data
                return Ok(None)

            case Err(error):
                match error:
                    case httpx.HTTPStatusError():
                        result = await self.create_user()
                        match result:
                            case Ok(user):
                                self.user = user
                                return Ok(None)

                            case Err(error):
                                return Err(error)

                    case httpx.ConnectError():
                        print(f"Bot #{self.id} - Network unreachable.")

                    case httpx.TimeoutException():
                        print(f"Bot #{self.id} - Request timed out.")

                    case _:
                        print(f"Bot #{self.id} - {type(error).__name__}: {error}")

                return Err(error)

    async def get_user(self) -> Result[User, Exception]:
        try:
            result = self.get_api_key()
            api_key: str | None = None
            match result:
                case Ok(data):
                    api_key = data
                case Err(error):
                    return Err(error)

            headers = {"X-API-KEY": api_key}
            response = await self.client.get("/users", headers=headers)

            response.raise_for_status()
            data = response.json()["data"]

            user = User.model_validate(data)
            return Ok(user)

        except Exception as e:
            return Err(e)

    async def create_user(self) -> Result[User, Exception]:
        try:
            result = self.get_api_key()
            api_key: str | None = None
            match result:
                case Ok(data):
                    api_key = data
                case Err(error):
                    return Err(error)

            headers = {"X-API-KEY": api_key}
            payload = {
                "first_name": "Portman",
                "last_name": f"Bot_{self.id}",
                "phone_number": f"+{self.id}",
                "email_address": f"portman_bot#{self.id}@system.local",
                "wallet_address": f"0x0000{self.id}",
                "role": "AUTOMATED",
                "password": f"bot#{self.id}_password",
            }

            response = await self.client.post("/users", headers=headers, json=payload)
            response.raise_for_status()

            data = response.json()["data"]
            user = User.model_validate(data)

            return Ok(user)

        except Exception as e:
            return Err(e)

    def get_api_key(self) -> Result[str, Exception]:
        try:
            email = f"portman_bot#{self.id}@system.local"
            secret = os.getenv("SYSTEM_SECRET")

            if secret is None:
                raise ValueError("SYSTEM_SECRET environment variable is not set")

            email_bytes = email.encode("utf-8")
            secret_bytes = secret.encode("utf-8")

            h = hmac.new(secret_bytes, email_bytes, digestmod=hashlib.sha256)
            return Ok(h.hexdigest())

        except Exception as e:
            return Err(e)
