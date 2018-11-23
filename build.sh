#!/bin/bash -eu

cd "$(dirname "${0}")"

# Support go 1 release 1.11 or higher.
GO_MAJOR_VERSION=1
MIN_GO_MINOR_VERSION=11

unset CGO_ENABLED
GO_VERSION="$(go version 2>/dev/null | cut -f3 -d' ')"
GO_MAJ="$(echo "${GO_VERSION}" | cut -f1 -d'.')"
GO_MIN="$(echo "${GO_VERSION}" | cut -f2 -d'.')"
if [ -z "${GO_VERSION}" ]; then
  echo "Have you installed go? I get no result from \`go version\` command." >&2
  exit 64
elif [ "${GO_MAJ}" != "go${GO_MAJOR_VERSION}" ] || [ "${GO_MIN}" -lt "${MIN_GO_MINOR_VERSION}" ]; then
  echo "Go version go${GO_MAJOR_VERSION}.x needed, where x >= ${MIN_GO_MINOR_VERSION}, but the version I found is: '${GO_VERSION}'" >&2
  echo "I found it here:" >&2
  which go >&2
  echo "The complete output of \`go version\` command is:" >&2
  go version >&2
  exit 65
fi
echo "Go version ok (${GO_VERSION} >= go${GO_MAJOR_VERSION}.${MIN_GO_MINOR_VERSION})"
TEST=false
OUTPUT_ALL_PLATFORMS="Building just for the platform in GOOS/GOARCH (${GOOS:-system default}/${GOARCH:-system default}) (build.sh -a argument NOT specified)"
OUTPUT_TEST="Test flag NOT detected (-t) as argument to build.sh script"
ALL_PLATFORMS=false
while getopts ":at" opt; do
    case "${opt}" in
        a)  ALL_PLATFORMS=true
            OUTPUT_ALL_PLATFORMS="Building for all platforms (build.sh -a argument specified)"
            ;;
        t)  TEST=true
            OUTPUT_TEST="Test flag detected (-t) as build.sh argument"
            ;;
    esac
done
echo "${OUTPUT_ALL_PLATFORMS}"
echo "${OUTPUT_TEST}"

export PATH="$(go env GOPATH)/bin:${PATH}"

function install {
  if [ "${1}" != 'native' ]; then
    local GOOS
    local GOARCH
    export GOOS="${1}"
    export GOARCH="${2}"
  fi
  CGO_ENABLED=0 go get -ldflags "-X main.revision=$(git rev-parse HEAD)" -v ./...
  go vet ./...

  # note, this just builds tests, it doesn't run them!
  go list ./... | while read PACKAGE; do
    CGO_ENABLED=0 go test -c "${PACKAGE}"
  done
}

if ${ALL_PLATFORMS}; then
  install windows 386
  install windows amd64
  install darwin  386
  install darwin  amd64
  install linux   386
  install linux   amd64
  install linux   arm
  install linux   arm64
else
  install native
fi

find "$(go env GOPATH)/bin" -name 'tc-client-generator*'

# capital X here ... we only want to delete things that are ignored!
git clean -fdX

if $TEST; then
  TEMP_SINGLE_REPORT="$(mktemp -t coverage.tmp.XXXXXXXXXX)"
  echo "mode: atomic" > coverage.report
  HEAD_REV="$(git rev-parse HEAD)"
  # Dump package list to file rather than pipe, to avoid exit inside loop not
  # causing outer shell to exit due to running in a subshell.
  PACKAGE_LIST="$(mktemp -t package-list.tmp.XXXXXXXXXX)"
  go list ./... > "${PACKAGE_LIST}"
  while read package
  do
    CGO_ENABLED=1 GORACE="history_size=7" go test -v -ldflags "-X github.com/taskcluster/tc-client-generator/cmd/tc-client-generator.revision=$(git rev-parse HEAD)" -race -timeout 1h -covermode=atomic "-coverprofile=${TEMP_SINGLE_REPORT}" "${package}"
    if [ -f "${TEMP_SINGLE_REPORT}" ]; then
      sed 1d "${TEMP_SINGLE_REPORT}" >> coverage.report
      rm "${TEMP_SINGLE_REPORT}"
    fi
  done < "${PACKAGE_LIST}"
  rm "${PACKAGE_LIST}"
fi

go vet ./...

GOOS= GOARCH= go get golang.org/x/lint/golint
golint ./...
GOOS= GOARCH= go get github.com/gordonklaus/ineffassign
ineffassign .

echo "Build successful!"
git status
