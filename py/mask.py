from typing import Any, Iterator


class Mask:

    def __init__(self) -> None:
        self._mask = 0
        self._count = 0

    def add(self, value: int) -> None:
        slot = _val_to_slot(value)
        self._set_mask(self._mask | slot)

    def pop(self, value: int) -> None:
        slot = _val_to_slot(value)
        return self._set_mask(self._mask & ~slot)

    def _set_mask(self, updated: int) -> None:
        if updated < self._mask:
            self._count -= 1
        elif updated > self._mask:
            self._count += 1

        self._mask = updated

    def __iter__(self) -> Iterator[int]:
        for val in range(1, 10):
            if _val_to_slot(val) & self._mask > 0:
                yield val

    def __contains__(self, o: Any) -> bool:
        if not isinstance(o, int):
            return False
        slot = _val_to_slot(o)
        return slot & self._mask > 0

    def __len__(self) -> int:
        return self._count


def _val_to_slot(val: int) -> int:
    return 1 << (val - 1)
