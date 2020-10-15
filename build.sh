#!/bin/bash

NAME="Go-Mirai-Client"
OUTPUT_DIR="output"

PLATFORMS="darwin/amd64"                         # amd64 only as of go1.5
PLATFORMS="$PLATFORMS windows/amd64 windows/386" # arm compilation not available for Windows
PLATFORMS="$PLATFORMS linux/amd64 linux/386"
PLATFORMS="$PLATFORMS linux/ppc64 linux/ppc64le"
PLATFORMS="$PLATFORMS linux/mips64 linux/mips64le" # experimental in go1.6
PLATFORMS="$PLATFORMS freebsd/amd64"
PLATFORMS="$PLATFORMS netbsd/amd64"          # amd64 only as of go1.6
PLATFORMS="$PLATFORMS openbsd/amd64"         # amd64 only as of go1.6
PLATFORMS="$PLATFORMS dragonfly/amd64"       # amd64 only as of go1.5
PLATFORMS="$PLATFORMS plan9/amd64 plan9/386" # as of go1.4
PLATFORMS="$PLATFORMS solaris/amd64"         # as of go1.3
PLATFORMS="$PLATFORMS linux/arm64"

for PLATFORM in $PLATFORMS; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  BIN_FILENAME="${OUTPUT_DIR}/${NAME}-${GOOS}-${GOARCH}"
  if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -v -ldflags \"-s -w\" -o ${BIN_FILENAME} $@"
  echo $CMD
  eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
done

# ARM builds
PLATFORMS_ARM="linux freebsd netbsd"
for GOOS in $PLATFORMS_ARM; do
  GOARCH="arm"
  # build for each ARM version
  for GOARM in 7 6 5; do
    BIN_FILENAME="${OUTPUT_DIR}/${NAME}-${GOOS}-${GOARCH}${GOARM}"
    CMD="GOARM=${GOARM} GOOS=${GOOS} GOARCH=${GOARCH} go build -v -ldflags \"-s -w\" -o ${BIN_FILENAME} $@"
    echo "${CMD}"
    eval "${CMD}" || FAILURES="${FAILURES} ${GOOS}/${GOARCH}${GOARM}"
  done
done
