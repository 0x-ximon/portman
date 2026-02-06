import asyncio
import httpx

from core.worker import Worker


class Manager:
    async def start(self, base_url: str, bots_amount: int) -> None:
        async with httpx.AsyncClient(
            base_url=base_url,
            headers={"Content-Type": "application/json"},
        ) as shared_client:
            bots = [Worker(i, shared_client) for i in range(1, bots_amount + 1)]
            await asyncio.gather(*(bot.run() for bot in bots))
