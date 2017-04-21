GO = go
ENV = export GOPATH=`pwd`/lib:`pwd`
GOGET = export GOPATH=`pwd`/lib; $(GO) get
BUILD = $(ENV); $(GO) build

all: stop_srv

dep:
	mkdir -p `pwd`/bin
	mkdir -p `pwd`/lib/src
	$(GOGET) golang.org/x/sys/windows/svc

stop_srv: dep
	$(BUILD) -o bin/StopSrv.exe src/main.go

clean:
	rm -rf ./bin/*
