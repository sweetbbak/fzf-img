build:
	go build -ldflags "-s -w"

test:
	fd . "$HOME/Pictures" --absolute-path -e png | head | ./fzfimg

install:
	/usr/bin/cp ./fzfimg "$HOME/bin"

compile:
	mkdir -p build
	mkdir -p build/mac
	env GOOS=darwin GOARCH=amd64 go build -o build/mac/fzfimg

	mkdir -p build/freebsd
	env GOOS=freebsd GOARCH=amd64 go build -o build/freebsd/fzfimg

	mkdir -p build/linux
	env GOOS=linux GOARCH=amd64 go build -o build/linux/fzfimg

	mkdir -p build/openbsd
	env GOOS=openbsd GOARCH=amd64 go build -o build/openbsd/fzfimg

	mkdir -p build/linux-arm
	env GOOS=linux GOARCH=arm go build -o build/linux-arm/fzfimg
