VERSION="Tetra-`git rev-parse --short HEAD`-`uname`-`uname -m`"

.PHONY: build clean run package docker-build docker-run test

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
	@mkdir build/doc
	@echo "Bundling documentation..."
	@godocdown . > build/doc/main.md
	@godocdown ./modules > build/doc/modules.md
	@godocdown ./atheme > build/doc/atheme.md
	@godocdown ./1459 > build/doc/1459.md
	@godocdown ./bot > build/doc/bot.md
	@cp -vrf doc/* build/doc/
	@mkdir build/etc
	@cp etc/config.yaml.example build/etc
	@cp README.md build
	@cp LICENSE build
	@echo "including source code"
	@mkdir build/src
	@cp -vrf ./1459 build/src
	@cp -vrf ./bot build/src
	@mv build ${VERSION}
	@tar cvzf ${VERSION}.tgz ${VERSION}
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

