import pytest
from unittest.mock import patch
from core.manager import Manager
from core.worker import Worker


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
