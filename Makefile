default: install build

frontend-install:
	@echo "** Installing Frontend dependencies via 'npm' **"
	@cd static && npm install

frontend-build:
	@echo "** Building Frontend **"
	@cd static && webpack

frontend-test:
	@echo "** Running tests for Frontend **"
	@echo "no test available yet for the Frontend"

backend-build:
	@echo "** Building Frontend **"
	go generate
	go build

backend-install:
	@echo "** Installing Backend dependencies via 'go get' **"
	@go get github.com/tools/godep
	@go get -u github.com/jteeuwen/go-bindata/...

backend-test:
	@echo "** Running tests for Backend  **"
	go test

install: frontend-install backend-install

build: frontend-build backend-build

test: frontend-test backend-test

release:
	gox -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"

