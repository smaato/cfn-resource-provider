.PHONY: build test clean

GOLANG_VERSION := 1.13.10

build:
	cfn generate
	docker run -it --workdir=/go/src/github.com/smaato/cfn-resource-provider/iam-samlprovider \
		-v $(PWD):/go/src/github.com/smaato/cfn-resource-provider/iam-samlprovider golang:$(GOLANG_VERSION) \
		go build -ldflags="-s -w" -tags="logging callback scheduler" -o bin/handler cmd/main.go

test:
	cfn generate
	docker run -it --workdir=/go/src/github.com/smaato/cfn-resource-provider/iam-samlprovider \
		-v $(PWD):/go/src/github.com/smaato/cfn-resource-provider/iam-samlprovider golang:$(GOLANG_VERSION) \
		env GOOS=linux go build -ldflags="-s -w" -o bin/handler cmd/main.go

clean:
	rm -rf bin
