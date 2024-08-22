.PHONY: build-example
build-example:
	GOOS=linux CGO_ENABLED=0 go build -o bin/example example/**
	zip bin/example.zip bin/example/bootstrap.zip
