.PHONY: \
	build \
	clean \
	run \
	_

export CGO_ENABLED=0

clean:
	@echo "Cleaning..."
	rm -rf ./build
	rm -rf ./web/build

build: web/build build/server

build/server:
	@echo "Building server..."
	mkdir -p build
	go build -ldflags '-extldflags "-static"' -tags timetzdata -o build/server ./cmd/server/main.go

web/node_modules:
	@echo "Installing node modules..."
	cd _ui && npm install

web/build: web/node_modules
	@echo "Building UI..."
	cd web && npm run build

run: web/build
	@echo "Running..."
	go run ./cmd/server/main.go