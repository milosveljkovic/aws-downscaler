#!/usr/bin/env bash

NAME="aws-downscaler"
VER="${1:-"v0.0.0"}"
BUILD_DATE=$(date +%Y-%m-%d-%H:%M:%S)
OUT_DIR="bins"

echo "VERSION: ${VER}"
echo "BUILD_DATE: ${BUILD_DATE}"

if [ -d "${OUT_DIR}" ]; then
    rm -rf "${OUT_DIR}"
fi
mkdir -p "${OUT_DIR}"

platforms=("windows/amd64" "windows/386" "darwin/amd64" "darwin/arm64" "linux/arm64" "linux/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="${NAME}-${GOOS}-${GOARCH}"

    echo "Building ${output_name} ..."

    env CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -installsuffix cgo \
        -buildvcs=false \
        -o "${OUT_DIR}/${output_name}" \
        -ldflags "\
            -X 'main.version=${VER}-${BUILD_DATE}' \
        " ./cmd/aws-downscaler

    if [ $? -ne 0 ]; then
        echo 'An error has occurred during 'GOOS=${GOOS} GOARCH=${GOARCH} go build'!'
        exit 1
    fi

    echo "Bin stored in ${OUT_DIR}/${output_name} ..."
    echo "---"
done