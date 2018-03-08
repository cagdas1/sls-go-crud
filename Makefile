build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/apis/create apis/create.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/apis/list apis/list.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/apis/get apis/get.go