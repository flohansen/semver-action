# SemVer GitHub Action

## Usage

```yaml
jobs:
  version:
    name: Determine Version
    runs-on: ubuntu-latest
    outputs:
      new-release: ${{ steps.semver.outputs.new-release }}
      new-release-version: ${{ steps.semver.outputs.new-release-version }}
    steps:
      - name: SemVer
        id: semver
        uses: flohansen/semver-action@v2
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
```
