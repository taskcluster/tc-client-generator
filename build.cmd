:: main script to build and test tc-client-generator on windows
@echo on
:: cd to parent dir of this script
pushd %~dp0
set CGO_ENABLED=0
go get -ldflags "-X main.revision=%REVISION%" -v -t ./... || exit /b 64
git rev-parse HEAD > revision.txt
set /p REVISION=< revision.txt
del revision.txt
go test -v -ldflags "-X github.com/taskcluster/tc-client-generator/cmd/tc-client-generator.revision=%REVISION%" ./... || exit /b 65
go get -v github.com/gordonklaus/ineffassign || exit /b 66
ineffassign . || exit /b 67
