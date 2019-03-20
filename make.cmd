@echo off
echo +-------------------------------------+
echo ^|              ^>^> httpd ^<^<           ^|
echo +-------------------------------------+

set GOPATH=
set GOROOT=

rmdir /Q /S bin>nul 2>nul

for /f %%i in ('chdir') do set pwd=%%i

for /f %%k in ('go env GOPATH') do set gp=%%k
for /f %%k in ('go env GOROOT') do set gr=%%k

set GOROOT=%gr%
set GOPATH=%gp%

echo step 1/3: create build output directory.
IF NOT EXIST bin mkdir bin

echo step 2/3: install libs...
go get github.com/urfave/cli

set GOPATH=%gp%;%pwd%

echo step 3/3: build...
go build -i -o bin/http.exe src/http.go
go build -i -o bin/proxy.exe src/proxy.go

echo build success!

