BINARY=engine
VERSION=$(shell cat version)

test: vendor
	go test -v -cover -covermode=atomic ./... 

check:
#	golangci-lint run
	gosec ./...
	staticcheck -f json ./...

vendor: 
	go mod vendor

#go get github.com/swaggo/swag/cmd/swag
docs-swag:
	swag init --parseDependency "github.com/gataca-io/vui-core/models"

#go get -u github.com/mikunalpha/goas
docs-oa:
	goas --module-path . --output docs/oas.json
	sed -i '' 's/\(github.*\.\)\(.*\)/\2/g' docs/oas.json 

docs: docs-swag docs-oa 

build: vendor
	go build -o ${BINARY}

unittest:
	go test -short $$(go list ./... | grep -v /vendor/)

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

tag: test
	git tag --delete v${VERSION} || true && git tag v${VERSION}

release: tag
	git push --delete origin v${VERSION} || true && git push origin v${VERSION}

.PHONY: test vendor build unittest clear tag release docs docs-oa docs-swag

