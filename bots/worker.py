import asyncio


class Worker:
    def __init__(self, bot_id):
        self.bot_id = bot_id

    async def run(self):
        print(f"Bot {self.bot_id} is starting...")
        await asyncio.sleep(1)
