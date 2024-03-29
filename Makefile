# GO_TEST_ARGS?=-v -count=1

all:
	go build ./...

test:
	go test $(GO_TEST_ARGS) ./pkg/...

itest:
	go test $(GO_TEST_ARGS) ./integration/...

clean:
	go clean ./...

.PHONY: all test clean
