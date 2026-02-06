import pytest

from unittest.mock import patch
from common.config import Config
from common.typings import Err


def test_get_args_custom_amount():
    """Test that custom amount flags are respected."""
    with patch("sys.argv", ["main.py", "--amount", "10"]):
        config = Config()

        result = config.load()
        if isinstance(result, Err):
            pytest.fail(f"Could not load config: {result.error}")

        assert config.bots_amount == 10
