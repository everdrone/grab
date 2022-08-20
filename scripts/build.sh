#!/usr/bin/env bash

if [ -z ${TARGET_GOOS+x} ]; then
  echo "- missing TARGET_GOOS env variable, setting to current os..."
  export TARGET_GOOS=$(go env GOOS)
fi

if [ -z ${TARGET_GOARCH+x} ]; then
  echo "- missing TARGET_GOARCH env variable, setting to current arch..."
  export TARGET_GOARCH=$(go env GOARCH)
fi

export BUILD_OUTPUT=$(go run ./scripts/get-artifact.go)
export COMMIT_SHA=$(git rev-parse HEAD)

export GOOS="${TARGET_GOOS}"
export GOARCH="${TARGET_GOARCH}"

export OPT_FLAGS="-s -w"

echo "- building..."

go build -o $BUILD_OUTPUT -ldflags="${OPT_FLAGS} -X 'github.com/everdrone/grab/internal/config.CommitHash=${COMMIT_SHA}'
            -X 'github.com/everdrone/grab/internal/config.BuildOS=${GOOS}' -X 'github.com/everdrone/grab/internal/config.BuildArch=${GOARCH}'"

echo "      os : ${GOOS}"
echo "    arch : ${GOARCH}"
echo "  output : ${BUILD_OUTPUT}"
echo "  commit : ${COMMIT_SHA}"
echo "   flags : ${OPT_FLAGS}"
