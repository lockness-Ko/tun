all: linux

windows: clean
	export CGO_ENABLED=0
	export GOOS=windows
	go build -ldflags="-extldflags=-static" -buildmode=pie .
	strip ./tun

linux: clean
	export CGO_ENABLED=0
	export GOOS=linux
	go build -ldflags="-extldflags=-static" -buildmode=pie .
	strip ./tun

package: build
	zip tun.zip ./tun

clean:
	-rm ./tun
	-rm ./tun.tar.gz