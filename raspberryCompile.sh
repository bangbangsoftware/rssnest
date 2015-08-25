#!/bin/bash
echo "Compiling for the raspberry pi"

echo "1. Set env vars"
export GOARCH=arm
export GOARM=5
export GOOS=linux
export GOPATH=/home/mick/work/rssnest

echo "2. Fixing imports"
goimports -w **/*.go
echo "3. Vetting"
go vet
echo "4. Building"
go build rssnest.go
echo "5. Testing (with coverage)"
go test -cover

echo "6. scp on to pi"
scp rssnest osmc@osmc:./rssnest/.
