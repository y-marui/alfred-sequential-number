"""Application orchestrator.

Note: Alfred invokes seq.py directly via inline scripts in info.plist.
This module provides programmatic access to the same logic for testing.
"""

from __future__ import annotations

from app.commands import seq_cmd


def generate(subcommand: str, args: list[str]) -> list[str]:
    """Generate sequential numbers.

    Args:
        subcommand: Format variant ('bin', 'oct', 'hex', 'Hex', 'alf', 'Alf', 'fmt', or '').
        args: Range/length tokens.

    Returns:
        List of generated strings.
    """
    return seq_cmd.generate(subcommand, args)
