"""Tests for seq_cmd."""

from __future__ import annotations

from app.commands.seq_cmd import generate


class TestDecimal:
    def test_length(self):
        assert generate("", ["5"]) == ["1", "2", "3", "4", "5"]

    def test_range(self):
        assert generate("", ["3-5"]) == ["3", "4", "5"]


class TestBinary:
    def test_length(self):
        assert generate("bin", ["4"]) == ["1", "10", "11", "100"]

    def test_range(self):
        assert generate("bin", ["3-4"]) == ["11", "100"]


class TestOctal:
    def test_length(self):
        assert generate("oct", ["8"]) == ["1", "2", "3", "4", "5", "6", "7", "10"]


class TestHexLower:
    def test_length(self):
        assert generate("hex", ["3"]) == ["1", "2", "3"]

    def test_includes_letters(self):
        result = generate("hex", ["16"])
        assert "a" in result
        assert "f" in result


class TestHexUpper:
    def test_includes_upper_letters(self):
        result = generate("Hex", ["16"])
        assert "A" in result
        assert "F" in result


class TestAlphaLower:
    def test_length(self):
        assert generate("alf", ["3"]) == ["a", "b", "c"]

    def test_wraps_to_aa(self):
        result = generate("alf", ["27"])
        assert result[-1] == "aa"


class TestAlphaUpper:
    def test_length(self):
        assert generate("Alf", ["3"]) == ["A", "B", "C"]

    def test_wraps_to_aa(self):
        result = generate("Alf", ["27"])
        assert result[-1] == "AA"


class TestFormat:
    def test_decimal_format(self):
        assert generate("fmt", ["%d", "3"]) == ["1", "2", "3"]

    def test_custom_prefix(self):
        result = generate("fmt", ["item-%d", "3"])
        assert result == ["item-1", "item-2", "item-3"]

    def test_multi_dimension(self):
        result = generate("fmt", ["#a-#", "3", "2"])
        assert result == ["a-1", "a-2", "b-1", "b-2", "c-1", "c-2"]

    def test_empty_args_returns_empty(self):
        assert generate("fmt", []) == []
