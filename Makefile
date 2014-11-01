VERSION="Tetra-0.1-`git rev-parse --short HEAD`-`uname`-`uname -m`"

.PHONY: build clean run package docker-build docker-run test pull debug

pull:
	git pull
	go get -v -u .

build:
	go build

debug:
	go clean
	go build -tags debug

clean:
	rm Tetra

run:
	@make build
	./Tetra

all:
	@make run

package:
	@echo "Building ${VERSION}..."
	@make build
	@echo "Setting up build prefix"
	@mkdir build
	@echo "Packing ${VERSION}..."
	@upx Tetra --ultra-brute --preserve-build-id
	@cp Tetra build
	@echo "Packing modules..."
	@cp -rf modules build
	@rm build/modules/doc.go
	@mkdir -p build/doc/go/bot
	@mkdir -p build/doc/go/external
	@mkdir -p build/lib
	@cp -rf lib build
	@echo "Bundling documentation..."
	@godoc github.com/coreos/go-etcd/etcd > build/doc/go/external/etcd
	@godoc code.google.com/p/go-uuid/uuid > build/doc/go/external/go-uuid
	@godoc github.com/sjkaliski/go-yo     > build/doc/go/external/go-yo
	@godoc github.com/stevedonovan/luar   > build/doc/go/external/luar
	@godoc github.com/kolo/xmlrpc         > build/doc/go/external/xmlrpc
	@godoc github.com/coreos/go-etcd      > build/doc/go/external/go-etcd
	@godoc github.com/codegangsta/cli     > build/doc/go/external/cli
	@godoc gopkg.in/yaml.v1               > build/doc/go/external/yaml
	@godoc .                              > build/doc/go/main
	@godoc ./modules                      > build/doc/go/modules
	@godoc ./atheme                       > build/doc/go/atheme
	@godoc ./1459                         > build/doc/go/1459
	@godoc ./bot                          > build/doc/go/tetra
	@godoc ./bot/modes                    > build/doc/go/bot/modes
	@godoc ./bot/web                      > build/doc/go/bot/web
	@cp -rf doc/* build/doc/
	@mkdir build/etc
	@cp etc/config.yaml.example build/etc
	@cp -rf etc/sendfile build/etc
	@cp README.md build
	@cp LICENSE build
	@echo "including source code"
	@echo "including help files"
	@cp -rf ./help build/
	@mkdir build/var
	@mv build ${VERSION}
	@tar czvf ${VERSION}.tgz ${VERSION}
	@rm -rf ${VERSION}
	@echo "Package at ${VERSION}.tgz"

docker-build:
	docker build -t xena/tetra .

docker-run:
	make -C ./test/testnet/ircd kill || true
	make -C ./test/testnet/ircd run
	docker run -dit --name tetra-etcd xena/etcd-minimal /etcd || true
	docker run --rm --link tetra-ircd:ircd --link tetra-etcd:etcd -it --name tetra xena/tetra .

test:
	make -C test test-build test-docker

