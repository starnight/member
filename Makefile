ARCH=amd64
OS=linux
CGO_FLAG=0
OUTPUT=webserver

all:
	CGO_ENABLE=${CGO_FLAG} GOOS=${OS} GOARCH=${ARCH} go build -o ${OUTPUT}

t := "/tmp/go-cover.$(shell /bin/bash -c "date +%Y%m%d%H%M%S").tmp"

test:
	go test -coverprofile=$t ./... && go tool cover -html=$t && unlink $t; \
	rm database/test.db

clean:
	rm ${OUTPUT}
