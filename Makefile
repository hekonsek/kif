PACKAGES := \
	github.com/hekonsek/kif/main \
	github.com/hekonsek/kif/main/cmd

all: format rice build silent-test

rice:
	(cd main/cmd && rice embed-go)

build:
	go build -o bin/kif main/kif.go

test:
	go test -v $(PACKAGES)

silent-test:
	go test $(PACKAGES)

format:
	go fmt $(PACKAGES)
