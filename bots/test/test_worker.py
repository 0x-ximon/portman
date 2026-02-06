import pytest
from core.worker import Worker


@pytest.mark.asyncio
async def test_worker_run(capsys):
    """Test if a single worker starts and prints its ID."""
    worker = Worker(bot_id=99)
    await worker.run()

    captured = capsys.readouterr()
    assert "Bot 99 is starting..." in captured.out
