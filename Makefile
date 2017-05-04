DEBUG_FLAG = $(if $(DEBUG),-debug)

deps:
	go get github.com/andrew-d/go-termutil/...
	go get github.com/urfave/cli/...
	go get github.com/mgutz/ansi/...
	go get -u github.com/golang/lint/golint
	go get -d -t ./...

build:
	go build

test: deps
	go test -v ./...

install: deps
	go install

fmt:
	go fmt ./...

lint:
	golint .

pkg: deps
	go get github.com/mitchellh/gox/...
	go get github.com/tcnksm/ghr
	mkdir -p pkg && cd pkg && gox ../...

clean:
	rm -f lltsv
