from unittest.mock import patch

from common.config import Config


def test_config_load(monkeypatch):
    monkeypatch.setenv("API_URL", "https://api.test.com")

    with patch("sys.argv", ["main.py", "--amount", "10"]):
        config = Config()
        config.load()

        assert config.bots_amount == 10
