import argparse


BOTS_AMOUNT = 50


def get_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Portman Bots CLI")
    parser.add_argument(
        "-a",
        "--amount",
        type=int,
        default=BOTS_AMOUNT,
        help="The amount of bots to create",
    )

    return parser.parse_args()
