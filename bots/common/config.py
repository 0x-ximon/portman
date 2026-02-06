import argparse
import os

from common.typings import Result, Ok, Err
from common.constants import BOTS_AMOUNT


class Config:
    env: os._Environ[str]
    args: argparse.Namespace

    def load(self) -> Result[None, Exception]:
        try:
            self.env = os.environ
            self.args = self.__get_args()

            return Ok(None)
        except Exception as e:
            return Err(e)

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
