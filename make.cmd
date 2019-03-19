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

echo step 1/2: create build output directory.
IF NOT EXIST bin mkdir bin

set GOPATH=%gp%;%pwd%

echo step 2/2: build...
go build -i -o bin/httpd.exe src/main.go

echo build success!

