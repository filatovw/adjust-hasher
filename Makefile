PHONY:build
build:
	go build -o ./hasher/bin/hasher ./hasher/cmd/cli

PHONY:clean
clean:
	rm -f ./hasher/bin/*

PHONY:test
test:
	go test -race -v ./...

PHONY:check
check:
	go vet ./...

PHONY:install
install:
	go install ./...
