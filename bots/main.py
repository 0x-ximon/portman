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

    args = config.args

    manager = Manager(args.amount)
    await manager.start()


if __name__ == "__main__":
    asyncio.run(main())
