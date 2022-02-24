
build:
	env CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o kratos-id-cleanup