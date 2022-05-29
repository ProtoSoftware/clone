<div align="center">

<img src="./.assets/icon.png" alt="Proto Logo" width="10%">
  
Fetching and clone release/tag data from one repository to another

</div>

### Note
Clone has been created with usage for Proto Curated repositories, you're free to use this with your own projects but keep in mind that no support will be offered and breaking changes could be introduced at any moment. 

## Usage
*It is recommended to use Clone with GitHub Actions schedules to keep repositories in sync on a schedule without needing to run it manually.*

Download a build from the release page for your system and run it in your command line. You **must** have `GITHUB_TOKEN`, `SOURCE_REPO` and `TARGET_REPO` In your environment beforehand, with `SOURCE_REPO` and `TARGET_REPO` being formatted like `ProtoSoftware/Example` instead of full URLs.

### GitHub Actions Example
```yml
name: Schedule

on:
  schedule:
    - cron: "0 0 * * *"
  workflow_dispatch:
  
jobs:
  Releases:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
        
      - name: Download Clone
        run: |
          DOWNLOAD_URL=$(curl -s https://api.github.com/repos/ProtoSoftware/clone/releases/latest \
          | grep browser_download_url \
          | grep linux_x86_64 \
          | cut -d '"' -f 4)
          curl -s -L --create-dirs -o ./clone "$DOWNLOAD_URL"
  
      - name: Make Bin Executable
        run: chmod +x ./clone

      - name: Check For Missing Releases
        run: GITHUB_TOKEN="${{ secrets.GITHUB_TOKEN }}" SOURCE_REPO="Example/Example" TARGET_REPO="Example/Example" ./clone
```

### Contributions
Please do not contribute to this repository unless you are a member of the organization, as noted before this is intended before internal use. However, feel free to make a fork of it and work with it yourself - following the [LICENSE](./LICENSE) of course.