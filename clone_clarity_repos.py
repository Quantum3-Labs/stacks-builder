import os
import subprocess
from tqdm import tqdm

# Folder for cloned repos
TARGET_DIR = "clarity_code_samples"

# List of GitHub repository URLs
repo_urls = [
    "https://github.com/clarity-lang/book.git",
    "https://github.com/clarity-lang/overview.git",
    "https://github.com/hirosystems/clarity-examples.git",
    "https://github.com/friedger/clarity-smart-contracts.git",
    "https://github.com/friedger/clarity-marketplace.git",
    "https://github.com/hirosystems/platform-template-nft-marketplace-dapp.git",
    "https://github.com/weavery/sworn.git",
    "https://github.com/CoinFabrik/stacy.git",
    "https://github.com/TheSoftNode/Crowd-Funding-App.git",
    "https://github.com/erfanyeganegi/droplinked-stacks-contract.git",
    "https://github.com/FLATLAY/droplinked-stacks-contract.git",
    "https://github.com/friedger/clarity-stacking-pools.git",
    "https://github.com/boomcrypto/clarity-deployed-contracts.git",
    "https://github.com/friedger/clarity-dao.git",
    "https://github.com/psq/swapr.git",
    "https://github.com/psq/flexr.git",
    "https://github.com/hirosystems/platform-template-nft-marketplace-dapp.git",
    "https://github.com/hirosystems/clarity-examples.git",
    "https://github.com/friedger/clarity-smart-contracts.git"
]

# Ensure the target folder exists
os.makedirs(TARGET_DIR, exist_ok=True)

for url in tqdm(repo_urls, desc="Processing Repositories", unit="repo"):
    repo_name = url.split("/")[-1].replace(".git", "")
    repo_path = os.path.join(TARGET_DIR, repo_name)

    # Skip if already cloned
    if os.path.exists(repo_path):
        continue

    # Try cloning
    try:
        subprocess.run(["git", "clone", "--depth", "1", url, repo_path], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    except subprocess.CalledProcessError:
        tqdm.write(f"X Failed to clone {url}")
        continue  # Skip to next repository

print("\nDone cloning all repositories.")