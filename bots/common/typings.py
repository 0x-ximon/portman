from typing import Literal, TypeVar

from pydantic import BaseModel, RootModel

from common.exceptions import Code

T = TypeVar("T")


class Error(BaseModel):
    code: Code
    detail: str


class Success[T](BaseModel):
    status: Literal["success"]
    data: T


class Failure(BaseModel):
    status: Literal["failure"]
    error: Error


Payload = RootModel[Success[T] | Failure]
