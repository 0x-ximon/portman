import asyncio

from common.config import Config
from common.typings import Err
from core.manager import Manager


async def main():
    config = Config()
    result = config.load()
    if isinstance(result, Err):
        print(f"Could not load config: {result.error}")
        return

    manager = Manager()
    await manager.start(config.api_url, config.bots_amount)


if __name__ == "__main__":
    asyncio.run(main())
