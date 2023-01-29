GO_PLUGIN_PATH := $(shell go env GOPATH)
PATH := $(GO_PLUGIN_PATH)/bin:$(PATH)

.PHONY: core test proto image clean fmt generate

core: proto generate
	go build -o build/core

test:
	go test -cpu 1,4 -timeout 7m github.com/SailGame/Core/...

proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/golang/mock/mockgen@v1.4.4
	mkdir -p pb
	protoc -I proto proto/core/*.proto --go_out=pb/ --go_opt=paths=source_relative --go-grpc_out=pb/ --go-grpc_opt=paths=source_relative
	cd pb/core && mockgen -destination=mocks/core.go -package=mocks . GameCore_ProviderServer,GameCore_ListenServer

image:
	docker build -t sailgame/core --build-arg GOPROXY=`go env GOPROXY` .

generate:
	go generate ./...

clean:
	rm -rf build/* pb/

fmt:
	go fmt github.com/SailGame/Core/...