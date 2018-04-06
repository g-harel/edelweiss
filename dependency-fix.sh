# moby client requires github.com/docker/go-connections/nat, but docker has vendored deps
rm -rf $GOPATH/src/github.com/docker/docker/vendor/github.com/docker/go-connections
go get -u github.com/pkg/errors
