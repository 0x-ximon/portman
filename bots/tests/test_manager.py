import pytest

from unittest.mock import patch
from core.manager import Manager
from core.worker import Worker


# @pytest.mark.asyncio
# async def test_manager_start_execution():
#     amount = 3
#     url = "http://localhost:3001"
#     manager = Manager()

#     with patch.object(Worker, "run", return_value=None) as mock_run:
#         await manager.start(url, amount)
#         assert mock_run.call_count == amount
