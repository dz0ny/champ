VERSION := 0.0.1

all: setup build

wercker: all test

setup:
	go get -v champ/...

clean:
	rm -f champ
	rm -rf pkg
	rm -rf bin
	find src/* -maxdepth 0 ! -name 'champ' -type d | xargs rm -rf

build:
	go build --ldflags '-w -X main.build=$(VERSION)' champ/cmd/champ
	go build --ldflags '-w -X main.build=$(VERSION)' champ/cmd/spinwheel

test:
	go get github.com/axw/gocov/gocov
	go get github.com/smartystreets/goconvey
	go test -v -race champ/...

rpi: setup
	go build -tags rpi --ldflags '-w -X main.build=$(VERSION)-rpi' champ/cmd/champ
	go build -tags rpi --ldflags '-w -X main.build=$(VERSION)-rpi' champ/cmd/spinwheel

cover:
	gocov test champ/... | gocov report
