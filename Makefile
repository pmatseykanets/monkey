default: test

clean:
	rm -f ./monkey

test:
	go vet ./...
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

tests: test

build: clean
	CGO_ENABLED=0 go build
