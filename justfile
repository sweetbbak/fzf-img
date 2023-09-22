build:
	go build -ldflags "-s -w"

test:
	fd . "$HOME/Pictures" --absolute-path -e png | head | ./fzfimg

install:
	/usr/bin/cp ./fzfimg "$HOME/bin"
