install:
	go build -o ~/bin/git-browse ./cmd/git-browse

dist:
	mkdir -p dist

build.all: build.linux build.darwin

build.linux: dist
	GOOS=linux go build -o dist/git-browse.linux ./cmd/git-browse

build.darwin: dist
	GOOS=darwin go build -o dist/git-browse.darwin ./cmd/git-browse
