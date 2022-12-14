app_name := proglog

CONFIG_PATH=$(HOME)/.proglog

TAG ?= 0.0.1

.PHONY: build
build :
	CGO_ENABLED=0 go build -o ./output/proglog -v ./cmd/proglog

build-linux :
	GOOS=linux CGO_ENABLED=0 go build -o ./output/linux/proglog -v ./cmd/proglog

.PHONY: run
run :
	$(default_exec)

.PHONY: proto
proto :
	protoc api/v1/*.proto --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --proto_path=.

.PHONY: version
version :
	git describe --tags > build/version
	git log --max-count=1 --pretty=format:%aI_%h > build/commit
	date +%FT%T%z > build/build_time

.PHONY: init
init :
	mkdir -p $(CONFIG_PATH)

.PHONY: gencert
gencert :
	cfssl gencert -initca test/ca-csr.json | cfssljson -bare ca
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=server test/server-csr.json | cfssljson -bare server
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=client -cn="root" test/client-csr.json | cfssljson -bare root-client
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=client -cn="nobody" test/client-csr.json | cfssljson -bare nobody-client
	mv *.pem *.csr ${CONFIG_PATH}

.PHONY: acl
acl :
	cp test/model.conf $(CONFIG_PATH)/model.conf
	cp test/policy.csv $(CONFIG_PATH)/policy.csv

build-docker:
	docker build -t github.com/neefrankie/proglog:$(TAG) .
