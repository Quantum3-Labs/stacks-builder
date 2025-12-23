#!/usr/bin/env python3
"""
Repository Cloning Script for Go Backend

This script clones Clarity repositories and reports progress to the Go backend.
Outputs newline-delimited JSON progress messages to stdout.
"""

import sys
import json
import subprocess
from pathlib import Path

# Get backend directory (1 level up from backend/scripts)
BACKEND_DIR = Path(__file__).parent.parent
TARGET_DIR = BACKEND_DIR / "data" / "clarity_code_samples"

REPO_URLS = [
    "https://github.com/hirosystems/clarity-examples.git",
    "https://github.com/hirosystems/platform-template-nft-marketplace-dapp.git",
    "https://github.com/friedger/clarity-marketplace.git",
    # "https://github.com/clarity-lang/book.git",
    # "https://github.com/clarity-lang/overview.git",
    # "https://github.com/friedger/clarity-smart-contracts.git",
    # "https://github.com/weavery/sworn.git",
    # "https://github.com/CoinFabrik/stacy.git",
    # "https://github.com/TheSoftNode/Crowd-Funding-App.git",
    # "https://github.com/erfanyeganegi/droplinked-stacks-contract.git",
    # "https://github.com/FLATLAY/droplinked-stacks-contract.git",
    # "https://github.com/friedger/clarity-stacking-pools.git",
    # "https://github.com/boomcrypto/clarity-deployed-contracts.git",
    # "https://github.com/friedger/clarity-dao.git",
    # "https://github.com/psq/swapr.git",
    # "https://github.com/psq/flexr.git",
]


def clone_repositories():
    """Clone all repositories with progress reporting"""
    # Ensure target directory exists
    TARGET_DIR.mkdir(parents=True, exist_ok=True)

    total = len(REPO_URLS)

    # Report start
    print(json.dumps({"type": "start", "total": total}), flush=True)

    cloned = 0
    skipped = 0
    failed = 0

    for i, url in enumerate(REPO_URLS, 1):
        repo_name = url.split("/")[-1].replace(".git", "")
        repo_path = TARGET_DIR / repo_name

        # Report progress
        print(json.dumps({
            "type": "progress",
            "current": i,
            "total": total,
            "message": f"Processing {repo_name}"
        }), flush=True)

        # Skip if already exists
        if repo_path.exists():
            skipped += 1
            continue

        # Try to clone
        try:
            subprocess.run(
                ["git", "clone", "--depth", "1", url, str(repo_path)],
                check=True,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
                timeout=60
            )
            cloned += 1
        except subprocess.TimeoutExpired:
            print(json.dumps({
                "type": "warning",
                "message": f"Timeout cloning {repo_name}"
            }), flush=True)
            failed += 1
        except subprocess.CalledProcessError:
            print(json.dumps({
                "type": "warning",
                "message": f"Failed to clone {repo_name}"
            }), flush=True)
            failed += 1
        except Exception as e:
            print(json.dumps({
                "type": "warning",
                "message": f"Error cloning {repo_name}: {str(e)}"
            }), flush=True)
            failed += 1

    # Report completion
    print(json.dumps({
        "type": "complete",
        "total_processed": total,
        "cloned": cloned,
        "skipped": skipped,
        "failed": failed
    }), flush=True)


if __name__ == "__main__":
    try:
        clone_repositories()
    except Exception as e:
        print(json.dumps({"type": "error", "message": str(e)}), file=sys.stderr)
        sys.exit(1)
