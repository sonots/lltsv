DEBUG_FLAG = $(if $(DEBUG),-debug)

build:
	GO111MODULE=on go build

test:
	GO111MODULE=on go test -v ./...

install:
	GO111MODULE=on go install

fmt:
	GO111MODULE=on go fmt ./...

lint:
	golint .

pkg:
	go get github.com/mitchellh/gox/...
	go get github.com/tcnksm/ghr
	mkdir -p pkg && cd pkg && gox ../...

clean:
	rm -f lltsv
