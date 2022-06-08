.PHONY: install
install:
	go build -o ~/bin/git-browse ./cmd/git-browse

dist:
	mkdir -p dist

.PHONY: clean
clean:
	rm -rf dist

.PHONY: build.all
build.all: build.linux build.darwin

.PHONY: build
build: dist
	go build -o dist/git-browse ./cmd/git-browse

.PHONY: build.linux
build.linux: dist
	GOOS=linux go build -o dist/git-browse.linux ./cmd/git-browse

.PHONY: build.darwin
build.darwin: dist
	GOOS=darwin go build -o dist/git-browse.darwin ./cmd/git-browse

