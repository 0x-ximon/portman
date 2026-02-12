import random
import asyncio
import httpx

from core.worker import Worker
from common.models import Ticker, Status


class Manager:
    tickers: list[Ticker] = []

    async def start(self, base_url: str, bots_amount: int) -> None:
        async with httpx.AsyncClient(
            base_url=base_url,
            headers={"Content-Type": "application/json"},
        ) as shared_client:
            await self.get_tickers(shared_client)

            bots = [
                Worker(i, shared_client, self.random_ticker())
                for i in range(1, bots_amount + 1)
            ]
            await asyncio.gather(*(bot.run() for bot in bots))

    # TODO: Implement error handling
    async def get_tickers(self, client: httpx.AsyncClient) -> None:
        response = await client.get("/tickers")
        response.raise_for_status()

        data = response.json()["data"]
        for t in data:
            ticker = Ticker.model_validate(t)
            if ticker.status == Status.OPEN:
                self.tickers.append(ticker)

    def random_ticker(self) -> str:
        return self.tickers[random.randint(0, len(self.tickers) - 1)].symbol
