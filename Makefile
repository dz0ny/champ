VERSION := 0.0.4

all: setup build

setup:
	go get -d -v champ/... # install all packages

dev-deps:
	go get -u github.com/axw/gocov/gocov
	go get -u github.com/laher/gols/cmd/...
	go get -u github.com/kardianos/govendor
	go get -u github.com/alecthomas/gometalinter
	bin/gometalinter --install --update
	go get -t -v champ/... # install test packages

clean:
	rm -f champ
	rm -rf pkg
	rm -rf bin
	find src/* -maxdepth 0 ! -name 'champ' -type d | xargs rm -rf
	find src/champ/vendor/* -maxdepth 0 ! -name 'champ' -type d | xargs rm -rf

build:
	go build --ldflags '-w -X main.build=$(VERSION)' champ/cmd/champ
	go build --ldflags '-w -X main.build=$(VERSION)' champ/cmd/spinwheel

rpi: setup
	go build -tags rpi --ldflags '-w -X main.build=$(VERSION)-rpi' champ/cmd/champ
	go build -tags rpi --ldflags '-w -X main.build=$(VERSION)-rpi' champ/cmd/spinwheel

install:
	install -m 755 champ /usr/sbin/
	install -m 755 spinwheel /usr/sbin/

lint:
	bin/gometalinter --fast --disable=gotype --disable=dupl --cyclo-over=30 --deadline=60s --exclude $(shell pwd)/src/champ/vendor src/champ/...
	find src/champ -not -path "./src/champ/vendor/*" -name '*.go' | xargs gofmt -w -s

test: lint cover
	go test -v -race $(shell go-ls champ/...)

cover:
	gocov test $(shell go-ls champ/...) | gocov report