.PHONY: build deps #release run test
	
# directory to output build
DIST_DIR=./dist
# get the date and time to use as a buildstamp
DATE=$$(date -Iseconds -u)
TIME=$$(date -Iseconds -u)
LDFLAGS="-s -w -X main.buildDate=$(DATE) -X main.buildTime=$(TIME)"

build:
	@go build --ldflags=$(LDFLAGS) -o $(DIST_DIR)/logvac main.go

deps:
	@go get -t -v ./...

docker:
	@docker build . -t thordin9/logvac