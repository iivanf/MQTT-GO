# Colors
NOCOLOR=\033[0m
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
FLAGVERSION=-ldflags "-s -w -X main.BuildTime=`date '+%Y-%m-%dT%H:%M:%SZ'` -X main.GitHash=`git rev-parse --short HEAD` -X main.GitTag=`git describe --exact-match --tags $(git log -n1 --pretty='%h') 2> /dev/null`"

all:
	@echo ''
	@echo '  make build      - Build station for x64 and boat for Raspberry'
	@echo ''

build: clean main copy-conf

clean:
	@echo "${BLUE}== Cleaning...${NOCOLOR}"
	rm -rf deploy/*
	@echo "${BLUE}== Cleaning...${NOCOLOR} ${GREEN} [OK] ${NOCOLOR}"

main:
	@echo "${BLUE}== Building station...${NOCOLOR}"
	go build ${FLAGVERSION} -v -o deploy/station/station cmd/station/station.go cmd/station/options.go
	@echo "${BLUE}== Building ...${NOCOLOR} ${GREEN} [OK] ${NOCOLOR}"
	@echo "${BLUE}== Building boat for Raspberry Pi ...${NOCOLOR}"
	CGO_ENABLED=1 GOARCH=arm GOARM=7 CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ LD=arm-linux-gnueabihf-ld go build ${FLAGVERSION} -v -o deploy/boat/boat cmd/boat/boat.go cmd/boat/options.go
	@echo "${BLUE}== Building for Raspberry Pi ...${NOCOLOR} ${GREEN} [OK] ${NOCOLOR}"


copy-conf: 
	@echo "${BLUE}== Copying conf...${NOCOLOR}"
	mkdir -p deploy/station/config
	mkdir -p deploy/boat/config
	cp cmd/station/config/config.toml deploy/station/config
	cp cmd/boat/config/config.toml deploy/boat/config
	@echo "${BLUE}== Copying conf...${NOCOLOR} ${GREEN} [OK] ${NOCOLOR}"


