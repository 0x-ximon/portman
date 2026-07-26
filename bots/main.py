import asyncio

from common.config import Config
from common.exceptions import PortmanException
from logic.manager import Manager


async def main():
    try:
        config = Config()
        config.load()

        manager = Manager()
        await manager.start(config.api_url, config.bots_amount)

    except PortmanException as e:
        print(f"Program crashed: {e}")


if __name__ == "__main__":
    asyncio.run(main())
