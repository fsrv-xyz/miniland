#!/usr/bin/env bash
set -e

declare PREFIX="VUE_APP_"

echo "${PREFIX}BUILD_INFO=$(git describe --tags --always --dirty)" > .env
echo "${PREFIX}PIPELINE_URL=${CI_PIPELINE_URL}" >> .env
echo "${PREFIX}SRC_REPO_LINK=${CI_PROJECT_URL}" >> .env