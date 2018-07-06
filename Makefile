PACKAGES := \
	github.com/hekonsek/skrt/main \
	github.com/hekonsek/skrt/main/cmd

all: format rice build silent-test

rice:
	(cd main/cmd && rice embed-go)

build:
	go build -o bin/skrt main/skrt.go

test:
	go test -v $(PACKAGES)

silent-test:
	go test $(PACKAGES)

format:
	go fmt $(PACKAGES)
