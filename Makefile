RELEASE_TAG=$$(git describe --abbrev=0 --tags)

dist-tools:
	@go get github.com/mitchellh/gox

dist: dist-tools
	rm -rf ./bin/*
	mkdir -p ./bin/iremail-cli_linux-amd64_$(RELEASE_TAG)
	mkdir -p ./bin/iremail-cli_linux-arm64_$(RELEASE_TAG)
	gox -osarch="linux/amd64" -output=./bin/iremail-cli_linux-amd64_$(RELEASE_TAG)/iremail-cli_$(RELEASE_TAG)
	gox -osarch="linux/arm64" -output=./bin/iremail-cli_linux-arm64_$(RELEASE_TAG)/iremail-cli_$(RELEASE_TAG)
	cd bin && ls --color=no | xargs -I {} tar -czf {}.tgz {}
	rm -rf ./bin/iremail-cli_linux-amd64_$(RELEASE_TAG)
	rm -rf ./bin/iremail-cli_linux-arm64_$(RELEASE_TAG)

release-tools:
	@go get github.com/tcnksm/ghr

release: release-tools
	ghr $(RELEASE_TAG) ./bin/

.PHONY: all dist-tools dist release-tools release
