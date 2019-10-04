# kp k+http
#  not to be confused with kadomtsev-petviashvili

all:
	# install go-1.13
	cp ../../k.go .
	echo -e "module github.com/ktye/i/_/kp\n\ngo 1.13" > go.mod
	GOPATH="" GO111MODULE=on go build -ldflags="-s -w"
	# there should now be an huge executable kp(.exe)
