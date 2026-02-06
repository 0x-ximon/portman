from dataclasses import dataclass
from typing import TypeVar, Generic, Union

T = TypeVar("T")
E = TypeVar("E")


@dataclass
class Ok(Generic[T]):
    data: T


@dataclass
class Err(Generic[E]):
    error: E


Result = Union[Ok[T], Err[E]]
