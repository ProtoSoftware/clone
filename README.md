# Clone
Tool used for fetching and mirroring releases from one repository to another

**Note:** This tool is intended to be used with Proto Curated repositories, you're free to use this with your own stuff but keep in mind it could have breaking changes or stop working completely at any moment. No support is offered for this tool.

## Usage
Download a build from the release page and run it in your command line, you must also have `GITHUB_TOKEN`, `SOURCE_REPO` and `TARGET_REPO` In your environment, with SOURCE_REPO and TARGET_REPO being formatted like `ProtoSoftware/Example` instead of full URLs.

If you would like to automate this process in your CI, as long as you make it auto download the latest build and give the binary executable permissions you should be good to go.

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
