test:
	go test ./... -v -race -cover -coverprofile=coverage.txt && go tool cover -func=coverage.txt

format:
	goimports -local "github.com/LasTshaMAN/txstore" -w ./
	# We need to run `gofmt` with `-s` flag as well (best practice, linters require it).
	# `goimports` doesn't support `-s` flag just yet.
	# For details see https://github.com/golang/go/issues/21476
	gofmt -w -s ./

lint:
	docker run --rm -v $(GOPATH)/pkg/mod:/go/pkg/mod:ro -v `pwd`:/`pwd`:ro -w /`pwd` golangci/golangci-lint:v1.27-alpine golangci-lint run --deadline=5m -v

build_console:
	go build -o ./bin/console ./cmd/
