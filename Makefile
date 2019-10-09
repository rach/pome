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
	@echo "** Building Backend **"
	go generate
	go build

backend-test:
	@echo "** Running tests for Backend  **"
	go test

install: frontend-install

build: frontend-build backend-build

test: frontend-test backend-test

release:
	gox -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"

