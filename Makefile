VERSION="Tetra-0.1-`git rev-parse --short HEAD`-`uname`-`uname -m`"

.PHONY: build clean run package docker-build docker-run test pull

pull:
	git pull
	go get -v -u .

build:
	go build

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
	@upx Tetra
	@cp Tetra build
	@echo "Packing modules..."
	@cp -rf modules build
	@rm build/modules/doc.go
	@mkdir -p build/doc/go/bot
	@mkdir -p build/doc/go/external
	@echo "Bundling documentation..."
	@godocdown github.com/coreos/go-etcd/etcd > build/doc/go/external/etcd.md
	@godocdown code.google.com/p/go-uuid/uuid > build/doc/go/external/go-uuid.md
	@godocdown github.com/codegangsta/negroni > build/doc/go/external/negroni.md
	@godocdown github.com/aarzilli/golua/lua  > build/doc/go/external/lua.md
	@godocdown github.com/drone/routes        > build/doc/go/external/routes.md
	@godocdown github.com/rcrowley/go-metrics > build/doc/go/external/go-metrics.md
	@godocdown github.com/sjkaliski/go-yo     > build/doc/go/external/go-yo.md
	@godocdown github.com/stevedonovan/luar   > build/doc/go/external/luar.md
	@godocdown gopkg.in/yaml.v1               > build/doc/go/external/yaml.md
	@godocdown .                              > build/doc/go/main.md
	@godocdown ./modules                      > build/doc/go/modules.md
	@godocdown ./atheme                       > build/doc/go/atheme.md
	@godocdown ./1459                         > build/doc/go/1459.md
	@godocdown ./bot                          > build/doc/go/tetra.md
	@godocdown ./bot/modes                    > build/doc/go/bot/modes.md
	@godocdown ./bot/web                      > build/doc/go/bot/web.md
	@cp -rf doc/* build/doc/
	@mkdir build/etc
	@cp etc/config.yaml.example build/etc
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
	docker run --rm --link tetra-ircd:ircd -it --name tetra xena/tetra .

test:
	make -C test test-build test-docker

