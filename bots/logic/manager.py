import asyncio
import random
from typing import Any, List

import httpx
import pydantic

from common.exceptions import ApiError, PortmanException, ValidationError
from common.logging import logger
from common.models import Ticker, TickerStatus
from common.typings import Failure, Payload, Success
from logic.worker import Worker


class Manager:
    tickers: list[Ticker] = []
    tasks: list[asyncio.Task] = []
    logger: Any

    def __init__(self) -> None:
        self.logger = logger.bind(id="manager")

    async def start(self, base_url: str, bots_amount: int) -> None:
        async with httpx.AsyncClient(
            base_url=base_url,
            headers={"Content-Type": "application/json"},
        ) as shared_client:
            self.logger.info("starting all bots", count=bots_amount)
            await self.get_tickers(shared_client)

            bots = [Worker(i, shared_client, self.random_ticker()) for i in range(1, bots_amount + 1)]
            self.tasks = [asyncio.create_task(bot.run()) for bot in bots]

            try:
                await asyncio.gather(*self.tasks)

            except (asyncio.CancelledError, KeyboardInterrupt):
                self.logger.info("shutdown signal received")

            finally:
                await self.stop()

    async def stop(self):
        self.logger.info("stopping all bots", count=len(self.tasks))
        for task in self.tasks:
            task.cancel()

        await asyncio.gather(*self.tasks, return_exceptions=True)
        self.logger.info("all bots stopped")

    def random_ticker(self) -> str:
        assert len(self.tickers) > 0, "no tickers available"
        return random.choice(self.tickers).symbol

    async def get_tickers(self, client: httpx.AsyncClient) -> None:
        try:
            response = await client.get("/tickers")
            payload = Payload[List[Ticker]].model_validate(response.json())

            match payload.root:
                case Success(data=ticker_list):
                    for ticker in ticker_list:
                        if ticker.status == TickerStatus.OPEN:
                            self.tickers.append(ticker)

                case Failure(error=err):
                    self.logger.error("failed to get tickers", code=err.code, detail=err.detail)
                    raise PortmanException.into(err.code, err.detail)

        except PortmanException:
            raise

        except httpx.RequestError as e:
            self.logger.error("request error occurred", detail=str(e))
            raise ApiError(f"Could not connect to the server: {e}") from e

        except pydantic.ValidationError as e:
            self.logger.error("failed during pydantic validation", detail=str(e))
            raise ValidationError("validation error occurred") from e

        except Exception as e:
            self.logger.error("an unexpected error occurred", detail=str(e))
            raise PortmanException(f"something went wrong: {e}") from e
