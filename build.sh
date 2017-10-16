#!/usr/bin/env bash

TARGETS=( darwin:amd64 linux:amd64 )

for target in ${TARGETS[@]}
do
    echo "Building ${target}..."
    export GOOS=$(echo ${target} | cut -d: -f1) GOARCH=$(echo ${target} | cut -d: -f2)
    go build -o "$(basename $(pwd))_${GOOS}_${GOARCH}"
done