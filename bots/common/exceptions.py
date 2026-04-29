from __future__ import annotations

from enum import Enum
from typing import Type


class Code(str, Enum):
    INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
    NOT_IMPLEMENTED = "NOT_IMPLEMENTED"
    UNAUTHORIZED = "UNAUTHORIZED"
    BAD_REQUEST = "BAD_REQUEST"
    UNEXPECTED = "UNEXPECTED"
    FORBIDDEN = "FORBIDDEN"
    NOT_FOUND = "NOT_FOUND"


class PortmanException(Exception):
    code: Code = Code.UNEXPECTED

    def __init__(self, detail: str) -> None:
        super().__init__(detail)

    @classmethod
    def into(cls, code: Code, detail: str) -> PortmanException:
        registry: dict[Code, Type[PortmanException]] = {
            Code.UNAUTHORIZED: AuthError,
            Code.FORBIDDEN: AuthError,
            Code.NOT_FOUND: NotFoundError,
            Code.BAD_REQUEST: ValidationError,
            Code.INTERNAL_SERVER_ERROR: ApiError,
            Code.NOT_IMPLEMENTED: ApiError,
        }

        error_class = registry.get(code, PortmanException)
        error = error_class(detail)
        error.code = code
        return error


class AuthError(PortmanException):
    code = Code.UNAUTHORIZED


class NotFoundError(PortmanException):
    code = Code.NOT_FOUND


class ValidationError(PortmanException):
    code = Code.BAD_REQUEST


class ApiError(PortmanException):
    code = Code.INTERNAL_SERVER_ERROR
