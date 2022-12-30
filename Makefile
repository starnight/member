ARCH=amd64
OS=linux
CGO_FLAG=0
OUTPUT=webserver

all:
	CGO_ENABLE=${CGO_FLAG} GOOS=${OS} GOARCH=${ARCH} go build -o ${OUTPUT}

test:
	go test

clean:
	rm ${OUTPUT}
