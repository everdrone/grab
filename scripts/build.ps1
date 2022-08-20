if ($null -eq $env:TARGET_GOOS) {
  echo "- missing TARGET_GOOS env variable, setting to current os..."
  $TARGET_GOOS = $(go env GOOS)
}

if ($null -eq $env:TARGET_GOARCH) {
  echo "- missing TARGET_GOARCH env variable, setting to current arch..."
  $TARGET_GOARCH = $(go env GOARCH)
}

$BUILD_OUTPUT = $(go run ./scripts/get-artifact.go)
$COMMIT_SHA = $(git rev-parse HEAD)

$GOOS = $TARGET_GOOS
$GOARCH = $TARGET_GOARCH

$OPT_FLAGS = "-s -w"

echo "- building..."

go build -o $BUILD_OUTPUT -ldflags="${OPT_FLAGS} -X 'github.com/everdrone/grab/internal/config.CommitHash=${COMMIT_SHA}'
            -X 'github.com/everdrone/grab/internal/config.BuildOS=${GOOS}' -X 'github.com/everdrone/grab/internal/config.BuildArch=${GOARCH}'"

echo "      os : ${GOOS}"
echo "    arch : ${GOARCH}"
echo "  output : ${BUILD_OUTPUT}"
echo "  commit : ${COMMIT_SHA}"
echo "   flags : ${OPT_FLAGS}"