#!/usr/bin/env python3

from __future__ import annotations

import json
import re
import subprocess
import sys


REPO = "laperezmu/fishing-duel"


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: github_issue_next_id.py <PREFIX>", file=sys.stderr)
        return 1

    prefix = sys.argv[1].strip().upper()
    output = subprocess.check_output(
        [
            "gh",
            "issue",
            "list",
            "--repo",
            REPO,
            "--state",
            "all",
            "--limit",
            "500",
            "--json",
            "title",
        ],
        text=True,
    )
    titles = [entry["title"] for entry in json.loads(output)]
    pattern = re.compile(rf"^\[{re.escape(prefix)}-(\d+)\]")
    seen = [int(match.group(1)) for title in titles if (match := pattern.match(title))]
    next_num = max(seen, default=0) + 1
    print(f"{prefix}-{next_num:02d}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
