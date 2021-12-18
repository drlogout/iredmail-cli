RELEASE_TAG=$$(git describe --abbrev=0 --tags)

dist-tools:
	@go install github.com/mitchellh/gox
	export PATH="$PATH:$HOME/go/bin"

dist: dist-tools
	rm -rf ./bin/*
	mkdir -p ./bin/iredmail-cli_linux-amd64_$(RELEASE_TAG)
	mkdir -p ./bin/iredmail-cli_linux-arm64_$(RELEASE_TAG)
	gox -osarch="linux/amd64" -output=./bin/iredmail-cli_linux-amd64_$(RELEASE_TAG)/iredmail-cli_$(RELEASE_TAG)
	gox -osarch="linux/arm64" -output=./bin/iredmail-cli_linux-arm64_$(RELEASE_TAG)/iredmail-cli_$(RELEASE_TAG)
	cd bin && ls --color=no | xargs -I {} tar -czf {}.tgz {}
	rm -rf ./bin/iredmail-cli_linux-amd64_$(RELEASE_TAG)
	rm -rf ./bin/iredmail-cli_linux-arm64_$(RELEASE_TAG)

release-tools:
	@go install github.com/tcnksm/ghr

release: release-tools
	ghr $(RELEASE_TAG) ./bin/

.PHONY: all dist-tools dist release-tools release
