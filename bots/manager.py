import asyncio
from worker import Worker


class Manager:
    bots: list[Worker]

    def __init__(self, amount: int):
        self.bots = [Worker(i) for i in range(amount)]

    async def start(self):
        await asyncio.gather(*(bot.run() for bot in self.bots))
