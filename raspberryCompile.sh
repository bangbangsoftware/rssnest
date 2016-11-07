#!/bin/bash
echo "If you haven't you should go download the latest go version (from https://golang.org/dl/)"
export GOPATH=$PWD:$PWD/config:$PWD/feeds
# apparently you don't set GOROOT export 
export GOROOT=/usr/share/go
set | grep GO
go get github.com/ChimeraCoder/anaconda
go get github.com/dutchcoders/goftp

cp *.go src/github.com/bangbangsoftware/.
cp -r feeds src/github.com/bangbangsoftware/.
cp -r config src/github.com/bangbangsoftware/.

go build FileStore
go build fileStore.go 
go build httpSource.go 
go build ftpTarget.go 
go build propergate.go 

echo "Compiling for the raspberry pi"

echo "1. Set env vars"
export GOARCH=arm
export GOARM=5
export GOOS=linux
# export GOPATH=/home/mick/work/rssnest

echo "2. Fixing imports"
./bin/goimports -w **/*.go
echo "3. Vetting"
go vet
echo "4. Building"
go build -o rssnest 
echo "5. Testing (with coverage)"
go test  -cover

echo "6. scp on to pi"
scp rssnest osmc@osmc:./rssnest/.
scp conf.json osmc@osmc:./rssnest/.
scp ./data/casts.json osmc@osmc:./rssnest/data/.
