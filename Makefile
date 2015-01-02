VERSION="Tetra-0.3-`git rev-parse --short HEAD`-`uname`-`uname -m`"

.PHONY: build clean run package docker-build docker-run test pull debug

pull:
	git pull
	godep go get -v .

build:
	godep go build

debug:
	godep go clean
	godep go build -tags debug

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
	@echo "Making debug binary"
	@make debug
	@upx Tetra --ultra-brute --preserve-build-id
	@cp Tetra build/Tetra-debug
	@echo "Packing modules..."
	@cp -rf modules build
	@rm build/modules/doc.go
	@mkdir -p build/doc/go/bot/script
	@mkdir -p build/lib
	@cp -rf lib build
	@echo "Bundling documentation..."
	@godoc .                              > build/doc/go/main
	@godoc ./modules                      > build/doc/go/modules
	@godoc ./atheme                       > build/doc/go/atheme
	@godoc ./1459                         > build/doc/go/1459
	@godoc ./bot                          > build/doc/go/tetra
	@godoc ./bot/modes                    > build/doc/go/bot/modes
	@godoc ./bot/web                      > build/doc/go/bot/web
	@godoc ./bot/script/crypto            > build/doc/go/bot/script/crypto
	@godoc ./bot/script/charybdis         > build/doc/go/bot/script/charybdis
	@godoc ./bot/script/strings           > build/doc/go/bot/script/strings
	@cp -rf doc/* build/doc/
	@mkdir build/etc
	@cp etc/config.yaml.example build/etc
	@cp -rf etc/sendfile build/etc
	@cp README.md build
	@cp LICENSE build
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

