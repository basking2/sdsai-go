# GO_TEST_ARGS?=-v -count=1

all:
	go build ./...

test:
	go test $(GO_TEST_ARGS) ./...

clean:
	go clean ./...

.PHONY: all test clean
