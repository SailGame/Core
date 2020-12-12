.PHONY: proto core image clean

core:
	go build -o build/core

proto:
	{\
    mkdir -p pb/; \
	protoc -I proto proto/core/core.proto --go-grpc_out=pb/; \
	}

image:
	docker build -t sailgame/core .

clean:
	rm -rf build/*