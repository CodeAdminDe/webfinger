#!/bin/bash
###################################################################
#  (c) 2025 Frederic Roggon                                       #
#                                                                 #
#  Licensed under the terms of GNU AFFERO GENERAL PUBLIC LICENSE. #
#  The full terms are provided via LICENSE file which is based    #
#  in the root of the code repository.                            #
#                                                                 #
#  Author: Frederic Roggon <frederic.roggon@codeadmin.de>         #
###################################################################
DATA_FILE=$RENOVATE_POST_UPGRADE_COMMAND_DATA_FILE
echo "Reading upgrade data from $DATA_FILE"

HAS_MAJOR=$(jq '[.[].updateType] | any(. == "major")' "$DATA_FILE")
HAS_MINOR=$(jq '[.[].updateType] | any(. == "minor")' "$DATA_FILE")
HAS_PATCH=$(jq '[.[].updateType] | any(. == "patch")' "$DATA_FILE")
HAS_DIGEST=$(jq '[.[].updateType] | any(. == "digest")' "$DATA_FILE")

if [ "$HAS_MAJOR" = "true" ]; then
    BUMP_TYPE="major"
elif [ "$HAS_MINOR" = "true" ]; then
    BUMP_TYPE="minor"
elif [ "$HAS_PATCH" = "true" ]; then
    BUMP_TYPE="patch"
elif [ "$HAS_DIGEST" = "true" ]; then
    BUMP_TYPE="digest"
fi

OLD_VERSION=$(cat .release-version | grep "version:" | awk '{print $2}')
IFS='.' read -r major minor patch <<< "$OLD_VERSION"

if [ "$BUMP_TYPE" = "major" ]; then
    major=$((major+1)); minor=0; patch=0;
elif [ "$BUMP_TYPE" = "minor" ]; then
    minor=$((minor+1)); patch=0
elif [ "$BUMP_TYPE" = "patch" ]; then
    patch=$((patch+1))
elif [ "$BUMP_TYPE" = "digest" ]; then
    patch=$((patch+1))
fi

CURRENT_VERSION="$major.$minor.$patch"

if [[ "$CURRENT_VERSION" =~ ^[0-9]+.[0-9]+.[0-9]+$ ]]; then
    #versionstring equals simple semver format. ok.
    echo "Old version in .release-version equals to $CURRENT_VERSION"
    echo "Renovate will bump to $CURRENT_VERSION after postUpgradeTask. Generating changelog entries..."
else
    echo "Failed to generate the current version string. The value does not equal the semver format."
    echo "Changelog not updated. Exit."
    exit 1
fi

DATE=$(date +'%Y-%m-%d')
UPGRADE_LINES=$(jq -r '.[] | "- \(.depName) \(.currentValue) => \(.newValue) (\(.updateType))"' "$DATA_FILE")
CHANGELOG_HEADER="## [v$CURRENT_VERSION] published at $DATE\n### Changes and/or dependency updates\n$UPGRADE_LINES\n"

if git tag -l | grep -q "^v$CURRENT_VERSION$"; then 
    echo "Version $CURRENT_VERSION already released. You need to update your release version to release."; exit 1;
else 
    echo "Version $CURRENT_VERSION is not released - validated successfully.";
fi

if grep -q "## \[$CURRENT_VERSION\]" CHANGELOG.md; then
    sed -i "/## \[$CURRENT_VERSION\]/,/\n## /c\\$CHANGELOG_HEADER" CHANGELOG.md
else
    echo -e "$CHANGELOG_HEADER\n\n$(cat CHANGELOG.md)" > CHANGELOG.md
fi

echo "Updated CHANGELOG.md for v$CURRENT_VERSION"
exit 0;