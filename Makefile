all: linux

build: windows linux

windows: clean
	bash -c "GOOS=windows go build ."
	strip ./tun.exe
	mv ./tun.exe ./build/tun.exe

linux: clean
	export CGO_ENABLED=0
	export GOOS=linux
	go build -ldflags="-extldflags=-static" -buildmode=pie .
	strip ./tun
	mv ./tun ./build/tun

package: build
	zip tun.zip ./build

clean:
	-rm ./build/tun
	-rm ./build/tun.exe
	-rm ./tun.tar.gz