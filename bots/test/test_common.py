import pytest
from common.config import Config
from common.typings import Err
from unittest.mock import patch


def test_get_args_custom_amount():
    """Test that custom amount flags are respected."""
    with patch("sys.argv", ["main.py", "--amount", "10"]):
        config = Config()
        result = config.load()

        if isinstance(result, Err):
            pytest.fail(f"Could not load config: {result.error}")

        args = config.args

        assert args.amount == 10
