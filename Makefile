


init:
	go install -v ./vendor/github.com/gogo/protobuf/protoc-gen-gogofast
	go install -v ./vendor/github.com/mwitkow/go-proto-validators/protoc-gen-govalidator

install:
	go install -v ./...
	go test -i ./...

test: install
	# go test -v -count 1 -race -short ./...
	go test -v -count 1 -race -timeout 30m ./tests
