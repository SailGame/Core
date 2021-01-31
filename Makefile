.PHONY: core test proto image clean fmt generate

core:
	go build -o build/core

test:
	go test -cpu 1,4 -timeout 7m github.com/SailGame/Core/...

proto:
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