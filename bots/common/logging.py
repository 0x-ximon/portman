from pathlib import Path

import structlog

# Ensure the directory exists
Path("tmp").mkdir(exist_ok=True)

structlog.configure(
    processors=[
        structlog.contextvars.merge_contextvars,
        structlog.processors.add_log_level,
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.dict_tracebacks,
        structlog.processors.JSONRenderer(),
    ],
    logger_factory=structlog.WriteLoggerFactory(file=open("tmp/dev.log", "a", encoding="utf-8")),
)

logger = structlog.get_logger()
