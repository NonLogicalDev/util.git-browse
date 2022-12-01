export prefix?=$(HOME)/.local
export bindir?=$(prefix)/bin

.PHONY: release
release:
	goreleaser release --rm-dist --snapshot --skip-publish

.PHONY: build
build:
	goreleaser build --rm-dist --snapshot --single-target --output dist/git-browse

.PHONY: install
install: build
	cp dist/frk $(bindir)
	
clean:
	rm -rf dist

