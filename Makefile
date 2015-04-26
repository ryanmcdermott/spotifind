# Prepend our vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
GOPATH := ${PWD}/vendor:${GOPATH}
export GOPATH

default: 
	build

build: 
	go build -v -o ./bin/spotifind ./src/spotifind

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt: 
	go fmt ./src/...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint: 
	golint ./src

run: build
	./bin/spotifind

test:
	go test ./src/...

vendor_clean:
	rm -dRf ./vendor/src

# We have to set GOPATH to just the vendor
# directory to ensure that `go get` doesn&#39;t
# update packages in our primary GOPATH instead.
# This will happen if you already have the package
# installed in GOPATH since `go get` will use
# that existing location as the destination.
vendor_get: vendor_clean
	GOPATH=${PWD}/vendor go get -d -u -v \
	github.com/rapito/go-spotify/spotify \
	github.com/codegangsta/cli

vendor_update: vendor_get
	rm -rf `find ./vendor/src -type d -name .git` \
    &amp;&amp; rm -rf `find ./vendor/src -type d -name .hg` \
    &amp;&amp; rm -rf `find ./vendor/src -type d -name .bzr` \
    &amp;&amp; rm -rf `find ./vendor/src -type d -name .svn`
