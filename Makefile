VERSION="Tetra-0.1-`git rev-parse --short HEAD`-`uname`-`uname -m`"

.PHONY: build clean run package docker-build docker-run test

build:
	go get -v -u .
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
	@mkdir -p build/doc/bot
	@echo "Bundling documentation..."
	@godocdown . > build/doc/main.md
	@godocdown ./modules > build/doc/modules.md
	@godocdown ./atheme > build/doc/atheme.md
	@godocdown ./1459 > build/doc/1459.md
	@godocdown ./bot > build/doc/bot.md
	@godocdown ./bot/modes > build/doc/bot/modes.md
	@godocdown ./bot/web > build/doc/bot/web.md
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

