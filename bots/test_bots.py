import pytest
import asyncio
from unittest.mock import patch, MagicMock
from cli import get_args
from worker import Worker
from manager import Manager


def test_get_args_custom_amount():
    """Test that custom amount flags are respected."""
    with patch("sys.argv", ["main.py", "--amount", "10"]):
        args = get_args()
        assert args.amount == 10


@pytest.mark.asyncio
async def test_worker_run(capsys):
    """Test if a single worker starts and prints its ID."""
    worker = Worker(bot_id=99)
    await worker.run()

    captured = capsys.readouterr()
    assert "Bot 99 is starting..." in captured.out


@pytest.mark.asyncio
async def test_manager_initialization():
    """Verify manager creates the correct number of worker instances."""
    amount = 5
    manager = Manager(amount)
    assert len(manager.bots) == amount
    assert isinstance(manager.bots[0], Worker)
    assert manager.bots[0].bot_id == 0


@pytest.mark.asyncio
async def test_manager_start_execution():
    """Verify manager triggers the run method for all bots."""
    amount = 3
    manager = Manager(amount)

    with patch.object(Worker, "run", return_value=None) as mock_run:
        await manager.start()
        assert mock_run.call_count == amount
