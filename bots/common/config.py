import argparse
import os

from dotenv import load_dotenv
from common.typings import Result, Ok, Err
from common.constants import BOTS_AMOUNT


class Config:
    api_url: str
    bots_amount: int

    def load(self) -> Result[None, Exception]:
        try:
            env = self.__get_env()
            args = self.__get_args()

            if args.amount is None:
                raise ValueError("Invalid amount")
            self.bots_amount = args.amount

            api_url = env.get("API_URL")
            if api_url is None:
                raise ValueError("API_URL    is not set")
            self.api_url = api_url

            return Ok(None)
        except Exception as e:
            return Err(e)

    def __get_env(self) -> os._Environ[str]:
        load_dotenv()
        return os.environ

    def __get_args(self) -> argparse.Namespace:
        parser = argparse.ArgumentParser(description="Portman Bots CLI")
        parser.add_argument(
            "-a",
            "--amount",
            type=int,
            default=BOTS_AMOUNT,
            help="The amount of bots to create",
        )

        return parser.parse_args()
