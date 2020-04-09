.PHONY: build test clean

build:
	cfn generate
	docker run -it --workdir=/go/src/github.com/smaato/cfn-resource-provider/iam-openidconnectprovider \
		-v $(PWD):/go/src/github.com/smaato/cfn-resource-provider/iam-openidconnectprovider golang:1.13.10 \
		go build -ldflags="-s -w" -tags="logging callback scheduler" -o bin/handler cmd/main.go

test:
	cfn generate
	env GOOS=linux go build -ldflags="-s -w" -o bin/handler cmd/main.go

clean:
	rm -rf bin
