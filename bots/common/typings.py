from typing import Generic, Literal, TypeVar, Union

from pydantic import BaseModel, RootModel

from common.exceptions import Code

T = TypeVar("T")


class Error(BaseModel):
    code: Code
    detail: str


class Success(BaseModel, Generic[T]):
    status: Literal["success"]
    data: T


class Failure(BaseModel):
    status: Literal["failure"]
    error: Error


Payload = RootModel[Union[Success[T], Failure]]
