.PHONY: build test clean

build:
	cfn generate
	docker run -it --workdir=/go/src/github.com/smaato/cfn-resource-provider/iam-samlprovider \
		-v $(PWD):/go/src/github.com/smaato/cfn-resource-provider/iam-samlprovider golang:1.12.14 \
		go build -ldflags="-s -w" -tags="logging callback metrics scheduler" -o bin/handler cmd/main.go

test:
	cfn generate
	env GOOS=linux go build -ldflags="-s -w" -o bin/handler cmd/main.go

clean:
	rm -rf bin