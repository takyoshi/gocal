GIT_VER:=$(shell git describe --tags)
DATE:=$(shell date +%Y-%m-%dT%H:%M:%SZ)

.PHONY: test get-deps install clean

all: test 

install:
	 cd cmd/gocal && go build -ldflags "-X=main.version ${GIT_VER} -X main.buildDate ${DATE}"
		install cmd/gocal/gocal ${GOPATH}/bin

get-deps:
	go get -t -d -v .
	cd cmd/gocal && go get -t -d -v .

packages:
	cd cmd/gocal && gox -os="linux darwin" -arch="amd64" -output "../../pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" -gcflags "-trimpath=${GOPATH}" -ldflags "-w -X main.version ${GIT_VER} -X main.buildDate ${DATE}"
	cd pkg && find . -name "*${GIT_VER}*" -type f -exec tar czf {}.tar.gz {} \;

test:
	go test -v

clean:
	rm -f cmd/gocal/gocal
	rm -f pkg/*

build:
	go build -gcflags="-trimpath=${HOME}" -ldflags="-w" cmd/gocal/gocal.go
