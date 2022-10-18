app_name := proglog

src_dir := $(current_dir)
out_dir := $(current_dir)/out

default_exec := $(out_dir)/$(app_name)

.PHONY: build
build :
	go build -o $(default_exec) -tags production -v $(src_dir)

.PHONY: run
run :
	$(default_exec)

.PHONY: proto
proto :
	protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.

.PHONY: version
version :
	git describe --tags > build/version
	git log --max-count=1 --pretty=format:%aI_%h > build/commit
	date +%FT%T%z > build/build_time
