export GOPATH=$(PWD)

HAS_STRIP := $(shell command -v strip 2> /dev/null)
HAS_UPX := $(shell command -v upx 2> /dev/null)

all: deps build compact

deps:
	go get -v -d luxaudio

build:
	go build -v luxaudio

compact:
ifdef HAS_STRIP
	strip luxaudio
endif
ifdef HAS_UPX
	upx -9 luxaudio
endif

clean:
	rm -rf src/github.com
	rm -rf bin
	rm -rf pkg
	rm -f luxaudio
