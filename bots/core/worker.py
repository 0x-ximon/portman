import asyncio
from common.typings import Result, Ok, Err


class Worker:
    def __init__(self, bot_id: int):
        self.bot_id = bot_id

    async def run(self):
        result = await self.connect()
        if isinstance(result, Err):
            print(f"Something went wrong: {result.error}")

        await asyncio.sleep(1)

    async def connect(self) -> Result[str, Exception]:
        try:
            print(f"Bot {self.bot_id} is starting...")
            await asyncio.sleep(1)
            return Ok("Connected")

        except Exception as e:
            return Err(e)
