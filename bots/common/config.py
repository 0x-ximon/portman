import argparse
import os

from dotenv import load_dotenv

from common.constants import BOTS_AMOUNT


class Config:
    api_url: str
    bots_amount: int

    def load(self) -> None:
        env = self.__get_env()
        args = self.__get_args()

        api_url = env.get("API_URL")
        if api_url is None:
            raise ValueError("API_URL is not set")

        if args.amount is None:
            raise ValueError("Amount is not set")

        self.api_url = api_url
        self.bots_amount = args.amount

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
