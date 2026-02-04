import asyncio
from cli import get_args
from manager import Manager


async def main():
    args = get_args()
    manager = Manager(args.amount)
    await manager.start()


if __name__ == "__main__":
    asyncio.run(main())
