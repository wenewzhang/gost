#!/bin/bash

enable_64k() {
	rm version.go
	cp version.lite version.go
	pushd ../../../x/internal/net
	rm transport.go
	cp transport.64 transport.go	
	popd  
}

enable_16k() {
	rm version.go
	cp version.mini version.go
	pushd ../../../x/internal/net
	rm transport.go
	cp transport.16 transport.go	
	popd 
}


if [ $# -eq 0 ];then
    echo "Usages: ./build.sh arm7/mips/mipsle/mips-mini/mipsle-mini/amd64/amd64-c/clean/log-on/log-off"
	exit
fi

if [ "$1" = 'arm7' ]; then   
	export GOOS=linux
	export GOARCH=arm
	export GOARM=7
	export CGO_ENABLED=0
	rm register.go
	cp register.client register.go
	enable_64k
	go build -o gost-arm7-client  -ldflags="-s -w"
	ls -la ./
fi


if [ "$1" = 'mips' ]; then   
	export CC=/home/jimmy/Downloads/openwrt-toolchain-23.05.0-rc3-ath79-generic_gcc-12.3.0_musl.Linux-x86_64/toolchain-mips_24kc_gcc-12.3.0_musl/bin/mips-openwrt-linux-musl-gcc
	export GOOS=linux
	export GOARCH=mips   
	export GOMIPS=softfloat
	export CGO_ENABLED=0
	rm register.go
	cp register.client register.go
	enable_64k
	go build -o gost-mips-client -ldflags="-s -w"
	ls -la ./
fi


if [ "$1" = 'mipsle' ]; then   
	export CC=/home/jimmy/Downloads/openwrt-toolchain-23.05.0-rc3-ath79-generic_gcc-12.3.0_musl.Linux-x86_64/toolchain-mips_24kc_gcc-12.3.0_musl/bin/mips-openwrt-linux-musl-gcc
	export GOOS=linux
	export GOARCH=mipsle   
	export GOMIPS=softfloat
	export CGO_ENABLED=0
	rm register.go
	cp register.client register.go
	enable_64k
	go build -o gost-mipsle-client -ldflags="-s -w"
	ls -la ./
fi

if [ "$1" = 'mipsle-mini' ]; then   
	export CC=/home/jimmy/Downloads/openwrt-toolchain-23.05.0-rc3-ath79-generic_gcc-12.3.0_musl.Linux-x86_64/toolchain-mips_24kc_gcc-12.3.0_musl/bin/mips-openwrt-linux-musl-gcc
	export GOOS=linux
	export GOARCH=mipsle   
	export GOMIPS=softfloat
	export CGO_ENABLED=0
	rm register.go
	cp register.client register.go
	enable_16k
	go build -o gost-mipsle-mini-client -ldflags="-s -w"
	ls -la ./
fi

if [ "$1" = 'mips-mini' ]; then   
	export CC=/home/jimmy/Downloads/openwrt-toolchain-23.05.0-rc3-ath79-generic_gcc-12.3.0_musl.Linux-x86_64/toolchain-mips_24kc_gcc-12.3.0_musl/bin/mips-openwrt-linux-musl-gcc
	export GOOS=linux
	export GOARCH=mips   
	export GOMIPS=softfloat
	export CGO_ENABLED=0
	rm register.go
	cp register.client register.go
	enable_16k
	go build -o gost-mips-mini-client -ldflags="-s -w"
	ls -la ./
fi

if [ "$1" = 'amd64' ]; then   
	export GOOS=linux
	export GOARCH=amd64   
	export CGO_ENABLED=0
	rm register.go
	cp register.server register.go
	enable_64k	
	go build -o gost-amd64-server -ldflags="-s -w"
	ls -la ./
fi

if [ "$1" = 'amd64-c' ]; then   
	export GOOS=linux
	export GOARCH=amd64   
	export CGO_ENABLED=0
	rm register.go
	cp register.client register.go
	enable_64k	
	go build -o gost-amd64-server -ldflags="-s -w"
	ls -la ./
fi

if [ "$1" = 'clean' ]; then   
	go clean -cache
	go clean -modcache
	go clean
fi

if [ "$1" = 'log-on' ]; then   
	cd ../../../x/logger
	rm logger.go
	cp logger.org logger.go
	echo "logger on!"
fi

if [ "$1" = 'log-off' ]; then   
	cd ../../../x/logger
	rm logger.go
	cp logger.nop logger.go
	echo "logger off!"
fi
