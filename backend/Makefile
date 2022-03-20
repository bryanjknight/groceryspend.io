.PHONY: build clean deploy

build: 
	go mod tidy
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/receipt/listRequests endpoints/receipt/list/list.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/receipt/createRequest endpoints/receipt/create/create.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

