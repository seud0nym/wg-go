VERSION ?= `[ -d ".git" ] && git describe --tags || date +%Y.%m.%d-dev`
LDFLAGS=-ldflags "-s -w -X main.appVersion=${VERSION}"
BINARY="wg-go"

build: *.go go.*
	go build ${LDFLAGS} -o ${BINARY}
	rm -rf /tmp/go-*

clean:
	rm -f ${BINARY}
