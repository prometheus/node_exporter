#!/usr/bin/env bash
set -euo pipefail

# Fetch the feature-gates manifest from remote. This is the source of truth.
FEATURE_GATES_MANIFEST_PATH="featuregate/feature-gates.md"
FEATURE_GATES_REMOTE_MANIFEST_RAW_URL="https://raw.githubusercontent.com/rexagod/node_exporter/feature-gated/$FEATURE_GATES_MANIFEST_PATH" # TODO: Point to prometheus remote.
FEATURE_GATES_REMOTE_MANIFEST_RAW_STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}\n" "$FEATURE_GATES_REMOTE_MANIFEST_RAW_URL")
if [ "$FEATURE_GATES_REMOTE_MANIFEST_RAW_STATUS_CODE" -ne 200 ]; then
  echo "Failed to fetch feature-gates data from remote, got $FEATURE_GATES_REMOTE_MANIFEST_RAW_STATUS_CODE."
  exit 1
fi
mkdir -p /tmp/featuregate
curl -s "$FEATURE_GATES_REMOTE_MANIFEST_RAW_URL" > "/tmp/$FEATURE_GATES_MANIFEST_PATH"

# Generate the updated feature-gates manifest locally. This is the working copy.
echo -e "<!-- This file is auto-generated. DO NOT EDIT MANUALLY. -->
# Feature gates
Below is the set of feature-gates currently present in the repository, along with the versions they were added and retired in.
| Name | Description | Adding Version | Retiring Version |
|------|-------------|----------------|------------------|" > "./$FEATURE_GATES_MANIFEST_PATH"
# For all non-test collector files, extract the feature-gates and their metadata.
find ./collector -type f ! -name "*_test.go" -exec cat {} \; | \
# * -0777: Slurp the entire file into a single string.
# * -n: Process the input line by line.
# * -e: Provide the script as an argument instead of reading it from a file.
# * $args =~ s/\n/ /g: Replace all newlines in the arguments with spaces (for multi-line arguments).
perl -0777 -ne '
    while (m/featuregate\.NewFeatureGate\(\s*([^)]*?)\s*\)/sg) {
        $args = $1;
        $args =~ s/\n/ /g;
        print "$args\n";
    }
' | \
# Squeeze multiple spaces into a single space.
tr -s ' ' | \
# Replace commas with pipes.
tr ',' '|' | \
# Remove double quotes.
sed 's/\"//g' >> "./$FEATURE_GATES_MANIFEST_PATH"

# Check for changes to existing feature-gates (instead of adding new ones) by comparing the local and remote feature-gates manifests.
if grep -Fxf "/tmp/$FEATURE_GATES_MANIFEST_PATH" "./$FEATURE_GATES_MANIFEST_PATH" > /dev/null; then
  exit 0
else
  echo "Manifest modifications are not consistent with the remote. Ensure that no existing feature-gate version commitments were changed."
fi
