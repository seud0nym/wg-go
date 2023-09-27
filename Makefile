APPVERSION ?= `[ -d ".git" ] && git describe --tags --long --dirty || date +%Y.%m.%d-dev`
WGCTRLVERSION ?= `grep wgctrl go.mod | cut -d' ' -f3`
LDFLAGS=-ldflags "-s -w -X main.appVersion=${APPVERSION} -X main.wgctrlVersion=${WGCTRLVERSION}"
BINARY="wg-go"

build: *.go go.*
	go build ${LDFLAGS} -o ${BINARY}

clean:
	rm -f ${BINARY}

arm:
	GOOS=linux GOARCH=arm GOARM=5 $(MAKE) build

update:
	go get -u
	go mod tidy
