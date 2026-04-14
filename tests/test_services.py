"""Tests for seq.py core functions."""

from __future__ import annotations

import sys
from pathlib import Path

# Ensure workflow/seq.py is importable.
_workflow_root = Path(__file__).resolve().parents[1] / "workflow"
if str(_workflow_root) not in sys.path:
    sys.path.insert(0, str(_workflow_root))

from seq import convert_Alphabet, convert_alphabet, load_range, loader  # noqa: E402


class TestConvertAlphabet:
    def test_1_is_a(self):
        assert convert_alphabet(1) == "a"

    def test_26_is_z(self):
        assert convert_alphabet(26) == "z"

    def test_27_is_aa(self):
        assert convert_alphabet(27) == "aa"

    def test_negative_raises(self):
        import pytest

        with pytest.raises(ValueError):
            convert_alphabet(0)


class TestConvertAlphabetUpper:
    def test_1_is_uppercase_a(self):
        assert convert_Alphabet(1) == "A"

    def test_26_is_uppercase_z(self):
        assert convert_Alphabet(26) == "Z"

    def test_27_is_uppercase_aa(self):
        assert convert_Alphabet(27) == "AA"


class TestLoadRange:
    def test_integer_input(self):
        assert load_range(3) == [1, 4]

    def test_digit_string(self):
        assert load_range("5") == [1, 6]

    def test_range_string(self):
        assert load_range("3-5") == [3, 6]


class TestLoader:
    def test_decimal(self):
        assert list(loader("%d", ["3"])) == ["1", "2", "3"]

    def test_binary(self):
        assert list(loader("%b", ["3"])) == ["1", "10", "11"]

    def test_hash_decimal(self):
        assert list(loader("#", ["3"])) == ["1", "2", "3"]

    def test_hash_alpha_lower(self):
        assert list(loader("#a", ["3"])) == ["a", "b", "c"]

    def test_hash_alpha_upper(self):
        assert list(loader("#A", ["3"])) == ["A", "B", "C"]

    def test_custom_format(self):
        result = list(loader("item-%d", ["3"]))
        assert result == ["item-1", "item-2", "item-3"]
